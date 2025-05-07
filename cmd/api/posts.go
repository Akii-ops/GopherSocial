package main

import (
	"backend/internal/store"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// for  ctx.WithValue
type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

// createPost godoc
//
//	@Summary		createPost
//	@Description	createPost
//	@Tags			createPost
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"CreatePostPayload"
//	@Success		201		{object}	store.Post
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Router			/posts/ [post]
//	@Security		ApiKeyAuth
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	authedUser := getUserFromContext(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  authedUser.ID,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// getPost godoc
//
//	@Summary		getPost
//	@Description	getPost
//	@Tags			getPost
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"getPost "
//	@Success		200		{object}	store.Post
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Router			/posts/{postID} [GET]
//	@Security		ApiKeyAuth
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return

	}

}

// deletePost godoc
//
//	@Summary		deletePost
//	@Description	deletePost
//	@Tags			deletePost
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"deletePost "
//	@Success		204		string		content
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Router			/posts/{postID} [delete]
//	@Security		ApiKeyAuth
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err.Error())
	}

	ctx := r.Context()

	err = app.store.Posts.Delete(ctx, int64(postID))

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponnse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// / 指针保证删除可行
type UpdatePostPayload struct {
	Content *string `json:"content" validate:"omitempty,max=1000"`
	Title   *string `json:"title" validate:"omitempty,max=100"`
}

// updatePost godoc
//
//	@Summary		updatePost
//	@Description	updatePost
//	@Tags			updatePost
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"updatePost "
//	@Param			payload	body		UpdatePostPayload	true	"UpdatePostPayload"
//	@Success		200		string		content
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Router			/posts/{postID} [patch]
//	@Security		ApiKeyAuth
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content

	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, int64(postID))

		if err != nil {

			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponnse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}

			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)

	return post
}
