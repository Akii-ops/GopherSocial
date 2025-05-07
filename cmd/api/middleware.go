package main

import (
	"backend/internal/store"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// read the auth header

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get base64

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}
			// /decode

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			// 实际username pwd
			username := app.config.auth.basic.user
			pwd := app.config.auth.basic.pwd
			//check
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pwd {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})

	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {

			//
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ") // authorization: Bearer <token>
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		// cache  -------------------
		ctx := r.Context()
		user, err := app.getUser(ctx, userID)
		if err != nil {
			///
			app.internalServerError(w, r, err)
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func (app *application) CheckPostOwnerShip(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := getUserFromContext(r)
		post := getPostFromCtx(r)

		// if it is the user's post
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check
		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			log.Printf("mddmmmd")
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r, fmt.Errorf("forbidden opration"))
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {

	role, err := app.store.Roles.GetByName(ctx, roleName)

	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil

}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {

	if !app.config.redis.enabled {
		return app.store.Users.GetByID(ctx, userID)
	}

	user, err := app.cachestore.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// set cache
		err = app.cachestore.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}

	}

	return user, nil

}

func (app *application) RateLimiterMiddlemare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})

}
