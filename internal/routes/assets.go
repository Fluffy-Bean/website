package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/data"
	"git.leggy.dev/Fluffy/Website/internal/web"
	"git.leggy.dev/Fluffy/Website/static"
)

func RegisterAssetsRoutes(h *web.Handler, r *chi.Mux) {
	r.Handle("/static/*", staticGet())
	r.Handle("/assets/*", assetsGet())
}

func staticGet() http.Handler {
	if _, err := os.Stat("static"); err != nil {
		return http.StripPrefix("/static/", http.FileServerFS(static.Dir))
	}
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}

func assetsGet() http.Handler {
	if _, err := os.Stat("data"); err != nil {
		return http.FileServerFS(data.Dir) // No need to strip prefix here as the folder is within /assets/*
	}
	return http.StripPrefix("/assets/", http.FileServer(http.Dir("data/assets")))
}
