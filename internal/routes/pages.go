package routes

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterPagesRoutes(h *web.Handler, r *chi.Mux) {
	r.Get("/", homeGet(h))
	r.Get("/fursona", fursonaGet(h))
}

func homeGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latestSong := h.LastFM.GetLatestSong()

		h.Template(w, r, "templates/pages/home.html", web.Data{
			"LatestSong": latestSong,
		})
	}
}

func fursonaGet(h *web.Handler) http.HandlerFunc {
	file, err := os.ReadFile(h.DataPath("art.json"))
	if err != nil {
		panic("read art file: " + err.Error())
	}

	var images []struct {
		File string `json:"file"`
	}
	if err := json.Unmarshal(file, &images); err != nil {
		panic("unmarshal art file: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "templates/pages/fursona.html", web.Data{
			"Art": images,
		})
	}
}
