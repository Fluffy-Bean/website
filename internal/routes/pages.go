package routes

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

type AristConfig struct {
	Name    string
	NSFW    bool
	Socials map[string]string
}

type ArtConfig struct {
	Path   string
	Artist string
	Sizes  map[int]string
}

type PagesConfig struct {
	GeneratedSizes []int
	Artists        []AristConfig
	Art            []ArtConfig
}

func RegisterPagesRoutes(h *web.Handler, r *chi.Mux, c PagesConfig) {
	r.Get("/", homeGet(h, &c))
	r.Get("/fursona", fursonaGet(h, &c))
}

func homeGet(h *web.Handler, c *PagesConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latestSong := h.LastFM.GetLatestSong()

		h.Template(w, r, "home.html", web.Data{
			"LatestSong": latestSong,
		})
	}
}

func fursonaGet(h *web.Handler, c *PagesConfig) http.HandlerFunc {
	for i := range c.Art {
		c.Art[i].Sizes = make(map[int]string)

		for _, size := range c.GeneratedSizes {
			src := strings.TrimPrefix(c.Art[i].Path, "/")
			sum := md5.Sum([]byte(src))
			str := fmt.Sprintf("/static/generated/%d/%x", size, sum)

			c.Art[i].Sizes[size] = str
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "fursona.html", web.Data{
			"Artists": c.Artists,
			"Art":     c.Art,
		})
	}
}
