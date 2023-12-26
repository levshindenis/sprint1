package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/middleware"
)

func MyRouter(hs handlers.HStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", middleware.WithCompression(middleware.WithLogging(hs.PostHandler)))
		r.Get("/{id}", middleware.WithCompression(middleware.WithLogging(hs.GetHandler)))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", middleware.WithCompression(hs.JSONPostHandler))
		})
	})
	return r
}
