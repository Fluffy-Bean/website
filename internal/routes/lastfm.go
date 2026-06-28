package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterLastFMRoutes(h *web.Handler, r *chi.Mux) {
	r.Handle("/last-fm/thumbnail", lastFMThumbnailGet(h))
}

func lastFMThumbnailGet(h *web.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastest := h.LastFM.GetLatestSong()

		if lastest == nil {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Write(lastest.Thumbnail)
	})
}
