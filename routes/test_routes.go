package routes

import (
	"net/http"

	"github.com/go-chi/chi"
)

func TestRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/internal-server-error", func(w http.ResponseWriter, r *http.Request) {
		panic("Forced internal server error")
	})

	return r
}
