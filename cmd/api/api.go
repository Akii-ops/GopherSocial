package main

import (
	"backend/internal/store"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr string

	maxIdleConns int
	maxOpenConns int
	maxIdleTime  string
}

// db
type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

		// /v1/health
		r.Get("/health", app.healthCheckHandler)

		// /v1/posts/
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {

				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)

				r.Delete("/", app.deletePostHandler)

				// 更新属性，put设置实体
				r.Patch("/", app.updatePostHandler)
			})

		})

	})

	return r
}

/*
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣶⣄⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣷⡀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⣿⣿⣿⣷⡄⠀
⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆
⠀⠀⠀⢀⣾⣿⡿⠟⠛⠉⠉⠉⠛⠻⢿⣿⣿⡄
⠀⢀⣴⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣦
⣰⣿⡟⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⢻⣿⣆
⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿
⣿⠀⠀⠀🐾⠀⠀⠀⠀⠀⠀⠀⠀🐾⠀⠀⠀⣿
⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿
⠸⣿⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⡿
⠀⠙⢿⣿⣦⣀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⣾⡿⠋
⠀⠀⠀⠈⠻⣿⣿⣷⣶⣶⣶⣶⣾⣿⣿⡿⠛⠁⠀⠀
⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀
*/

func (app *application) run(mux *chi.Mux) error {

	srv := &http.Server{

		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
