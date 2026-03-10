package main

import (
	"log"
	"net/http"
	"time"

	"github.com/farrasnazhif/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr string
	// conns = connections
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string // time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// POST /v1/posts
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			// GET, DELETE, UPDATE /v1/posts/postID
			r.Route("/{postID}", func(r chi.Router) {
				// to receive current post context (data) / fetching the post to receive current context (data)
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has starter at %s", app.config.addr)

	return srv.ListenAndServe()
}
