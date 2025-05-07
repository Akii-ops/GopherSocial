package main

import (
	"backend/internal/store"
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// getUser 处理获取用户信息的请求
//
//	@Summary		获取用户信息
//	@Description	根据用户ID返回用户详情
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int			true	"user ID"
//	@Success		200		{object}	store.User	"用户信息"
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [get]
//	@Security		ApiKeyAuth
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	//
	user, err := app.getUser(ctx, userID)
	// user, err := app.store.Users.GetByID(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponnse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return

	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// followUser godoc
//
//	@Summary		关注用户
//	@Description	关注用户，路径参数为被关注用户，
//	@Tags			users follow
//
//	@Param			userID	path	int	true	"the ID of user to be followed"
//	@Accept			json
//	@Success		204	{string}	string	"No Content"
//	@Failure		500	{object}	error
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Router			/users/{userID}/follow [put]
//	@Security		ApiKeyAuth
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, followedID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// no content
	if err := NoContentResponse(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// unfollowUser godoc
//
//	@Summary		取关用户
//	@Description	取消关注用户，路径参数为被取关用户，
//	@Tags			users unfollow
//	@Accept			json
//	@Param			userID	path		int		true	"the ID of user to be unfollowed"
//	@Success		204		{string}	string	"No Content"
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Router			/users/{userID}/unfollow [put]
//	@Security		ApiKeyAuth
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	unfollowedUser := getUserFromContext(r)

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	// TODO:
	if err := app.store.Followers.UnFollow(r.Context(), unfollowedUser.ID, unfollowedID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// no content
	if err := NoContentResponse(w); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// activateUser godoc
//
//	@Summary		激活用户
//	@Description	激活关注用户
//	@Tags			users activate
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	" "
//	@Success		204		{string}	string	"User activated"
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponnse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := NoContentResponse(w); err != nil {
		app.internalServerError(w, r, err)

	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponnse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return

		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)

	return user
}
