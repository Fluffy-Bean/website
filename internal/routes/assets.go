package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
	"git.leggy.dev/Fluffy/Website/static"
)

func RegisterAssetsRoutes(h *web.Handler, r *chi.Mux) {
	r.Handle("/static/*", staticGet())
}

func staticGet() http.Handler {
	if _, err := os.Stat("static"); err != nil {
		return http.StripPrefix("/static/", http.FileServerFS(static.Dir))
	}
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}
