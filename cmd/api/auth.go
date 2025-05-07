package main

import (
	"backend/internal/mailer"
	"backend/internal/store"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	User  *store.User
	Token string `json:"token"`
}

// userRegister godoc
//
//	@Summary		注册用户
//	@Description	注册用户
//	@Tags			users unfollow
//	@Accept			json
//	@Produce		json

// @Param		payload	body		RegisterUserPayload	true	"User credentials"
// @Success	200		{object}	UserWithToken
// @Failure	500		{object}	error
// @Failure	400		{object}	error
// @Router		/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// hashing password

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	//store user

	// token  hash
	plainToken := uuid.New().String()

	// store
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	// send email

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	// isDevEnv := app.config.env == "development"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	// TODO : 异步
	err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars)
	if err != nil {

		app.logger.Errorw("error sending welcome email", "error", err)

		//rollback user creation if email fails (SAGA pattern)

		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)

		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// email / [username] [phone number]  + password
type CreateTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72,min=3"`
}

// tokenGenerate godoc
//
//	@Summary		token
//	@Description	token
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTokenPayload	true	"User credentials"/
//	@Success		201		{object}	string
//	@Failure		500		{object}	error
//	@Failure		401		{object}	error
//	@Failure		400		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parse payload``

	var payload CreateTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fetch user
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			// app.notFoundResponnse(w, r, err)   枚举攻击
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	//generate token -> add claims

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send it to client

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}
