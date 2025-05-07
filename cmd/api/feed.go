package main

import (
	"backend/internal/store"
	"net/http"
)

// feed godoc
//
//	@Summary		feed
//	@Description	token
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//
//	@Param			limit	query		string	true	"limit"
//	@Param			offset	query		string	true	"limit"
//	@Param			sort	query		string	true	"limit"
//	@Param			tags	query		string	true	"limit"
//	@Param			search	query		string	true	"limit"
//	@Success		200		{object}	store.PostWithMetadata
//	@Failure		500		{object}	error
//	@Failure		400		{object}	error
//	@Router			/users/feed [get]
//
//	@Security		ApiKeyAuth
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// pagination, filters

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(114), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
