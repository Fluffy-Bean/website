package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"

	"git.leggy.dev/Fluffy/Website/internal/blog"
	"git.leggy.dev/Fluffy/Website/internal/web"
)

var blogs map[string]blog.Blog

func RegisterBlogRoutes(h *web.Handler, r *chi.Mux) {
	blogs = make(map[string]blog.Blog)

	files, err := os.ReadDir(h.DataPath("blogs"))
	if err != nil {
		panic("read blogs directory: " + err.Error())
	}

	for _, file := range files {
		var blogData blog.Blog

		blogData.Title = strings.TrimSuffix(file.Name(), ".md")

		f, err := os.ReadFile(h.DataPath("blogs/" + file.Name()))
		if err != nil {
			panic("read blog data: " + err.Error())
		}

		err = goldmark.Convert(f, &blogData.Data)
		if err != nil {
			panic("convert blog post: " + err.Error())
		}

		blogs[blogData.Title] = blogData
	}

	r.Get("/blogs", blogListGet(h))
	r.Get("/blogs/{blogID}", blogGet(h))
}

func blogListGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "templates/pages/blog_list.html", web.Data{
			"Blogs": blogs,
		})
	}
}

func blogGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blogID := chi.URLParam(r, "blogID")

		b, ok := blogs[blogID]
		if !ok {
			http.NotFound(w, r)

			return
		}

		h.Template(w, r, "templates/pages/blog_post.html", web.Data{
			"BlogHTML": b.Data.String(),
		})
	}
}
