package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/handler"
)

func RegisterAssetsRoutes(h *handler.Handler, r *chi.Mux) {
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(h.DataPath("assets")))))
}
