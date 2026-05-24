package routes

import (
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

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
		if strings.HasPrefix(file.Name(), "_") {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".md")

		parts := strings.Split(slug, "-")
		if len(parts) != 2 {
			panic("unexpected file name, want yyyy_mm_dd-snake_case_title: " + file.Name())
		}

		publishedAt, err := time.Parse("2006_01_02", parts[0])
		if err != nil {
			panic("unexpected file name: " + err.Error())
		}

		title := strings.ReplaceAll(parts[1], "_", " ")

		var blogData blog.Blog
		blogData.Slug = slug
		blogData.Title = title
		blogData.PublishedAt = publishedAt.UTC()

		f, err := os.ReadFile(h.DataPath("blogs/" + file.Name()))
		if err != nil {
			panic("read blog data: " + err.Error())
		}

		err = goldmark.Convert(f, &blogData.Data)
		if err != nil {
			panic("convert blog post: " + err.Error())
		}

		blogs[slug] = blogData
	}

	r.Get("/blogs", blogListGet(h))
	r.Get("/blogs/{blogID}", blogGet(h))
}

func blogListGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sorted []blog.Blog
		for _, b := range blogs {
			sorted = append(sorted, b)
		}

		slices.SortFunc(sorted, func(a, b blog.Blog) int {
			if !a.PublishedAt.Before(b.PublishedAt) {
				return -1
			}
			return 0
		})

		h.Template(w, r, "templates/pages/blog_list.html", web.Data{
			"Blogs": sorted,
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

		oldBlogTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		isOldBlog := b.PublishedAt.Before(oldBlogTime)

		h.Template(w, r, "templates/pages/blog_post.html", web.Data{
			"IsOldBlog": isOldBlog,
			"BlogHTML":  b.Data.String(),
		})
	}
}
