package main

import (
	"net/http"

	"github.com/farrasnazhif/go-social/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Get feed of posts from followed users with pagination, filters, and search
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Number of posts to return (1-20)"		minimum(1)	maximum(20)
//	@Param			offset	query		int		false	"Number of posts to skip"				minimum(0)
//	@Param			sort	query		string	false	"Sort order (asc or desc)"				Enums(asc, desc)
//	@Param			tags	query		string	false	"Filter by tags (comma separated, e.g. golang,backend)"
//	@Param			search	query		string	false	"Search in title or content"
//	@Param			since	query		string	false	"Filter posts after date (YYYY-MM-DD HH:MM:SS)"
//	@Param			until	query		string	false	"Filter posts before date (YYYY-MM-DD HH:MM:SS)"
//	@Success		200		{array}	store.PostWithMetadata
//	@Failure		400		{object}	error	"Invalid query parameters"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort
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

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(2), fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "Feed received successfully", feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
