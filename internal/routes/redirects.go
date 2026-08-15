package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterRedirectsRoutes(h *web.Handler, r *chi.Mux) {
	r.Get("/refsheet", redirectRoute("/fursona"))

	r.Get("/search/*", redirectRoute("/blogs"))

	r.Get("/posts", redirectRoute("/blogs"))
	r.Get("/posts/2024_11_07-umami", redirectRoute("/blogs/2024_11_07-Umami"))
	r.Get("/posts/2023_06_19-hello-django", redirectRoute("/blogs/2023_06_19-Hello_Django"))
	r.Get("/posts/2024_06_06-pebble-time", redirectRoute("/blogs/2024_06_06-Pebble_Time"))
	r.Get("/posts/2024_05_28-astro-is-hard", redirectRoute("/blogs/2024_05_28-Astro_is_hard"))
}

func redirectRoute(dist string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, dist, http.StatusMovedPermanently);
	}
}
