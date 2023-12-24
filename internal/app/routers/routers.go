package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/logging"
)

func MyRouter(hs handlers.HStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", logging.WithLogging(hs.PostHandler))
		r.Get("/{id}", logging.WithLogging(hs.GetHandler))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", hs.JSONPostHandler)
		})
	})
	return r
}
