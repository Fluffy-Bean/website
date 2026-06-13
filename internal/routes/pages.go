package routes

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

// These numbers come from the Makefile
var generatedSizes = []int{256, 512, 1024}

func RegisterPagesRoutes(h *web.Handler, r *chi.Mux) {
	r.Get("/", homeGet(h))
	r.Get("/fursona", fursonaGet(h))
}

func homeGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latestSong := h.LastFM.GetLatestSong()

		h.Template(w, r, "home.html", web.Data{
			"LatestSong": latestSong,
		})
	}
}

func fursonaGet(h *web.Handler) http.HandlerFunc {
	file, err := h.ReadDataFile("art.json")
	if err != nil {
		panic("read art file: " + err.Error())
	}

	var data struct {
		Artists map[string]struct {
			NSFW    bool              `json:"nsfw"`
			Socials map[string]string `json:"socials"`
		} `json:"artists"`
		Files []struct {
			Sizes  map[int]string
			Path   string `json:"path"`
			Artist string `json:"artist"`
		}
	}

	if err := json.Unmarshal(file, &data); err != nil {
		panic("unmarshal art file: " + err.Error())
	}

	for i := range data.Files {
		data.Files[i].Sizes = make(map[int]string)

		for _, size := range generatedSizes {
			src := strings.TrimPrefix(data.Files[i].Path, "/")
			sum := md5.Sum([]byte(src))
			str := fmt.Sprintf("/static/generated/%d/%x", size, sum)

			data.Files[i].Sizes[size] = str
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "fursona.html", web.Data{
			"Artists": data.Artists,
			"Files":   data.Files,
		})
	}
}
