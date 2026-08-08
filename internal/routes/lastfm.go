package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

var transparentGIF, _ = base64.RawStdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAACXBIWXMAAC4jAAAuIwF4pT92AAAAC0lEQVQI12NgAAIAAAUAAeImBZsAAAAASUVORK5CYII=")

func RegisterLastFMRoutes(h *web.Handler, r *chi.Mux) {
	r.Handle("/last-fm/thumbnail", lastFMThumbnailGet(h))
}

func lastFMThumbnailGet(h *web.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		latest := h.LastFM.GetLatestSong()

		if latest == nil {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if latest.Thumbnail == nil {
			w.Header().Set("Content-Type", "image/png")
			w.Write(transparentGIF)

			return
		}

		w.Write(latest.Thumbnail)
	})
}
