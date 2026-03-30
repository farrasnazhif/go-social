package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/farrasnazhif/go-social/docs" // this is required to generate swagger docs
	"github.com/farrasnazhif/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
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
		// GET /api/v1/health
		r.Get("/health", app.healthCheckHandler)

		// GET /api/v1/swagger
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// POST /api/v1/posts
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			// GET, DELETE, UPDATE /api/v1/posts/{postID}
			r.Route("/{postID}", func(r chi.Router) {
				// fetch post by ID, use for patch/put method
				// we need to fetch by ID first before updating
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		// GET /api/v1/users/{userID}
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// public routes

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Version = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has starter at", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
