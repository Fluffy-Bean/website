package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"git.leggy.dev/Fluffy/Website/internal/handler"
	"git.leggy.dev/Fluffy/Website/internal/lastfm"
	"git.leggy.dev/Fluffy/Website/internal/routes"
	"git.leggy.dev/Fluffy/Website/internal/sse"
)

const maxBodySize = 1 * 1024 * 1024 // 1mb

type flags struct {
	secret   string
	lastfm   string
	address  string
	port     int
	dataPath string
}

func main() {
	f := parseFlags()

	if f.secret == "" {
		panic("No secret provided")
	}
	if f.lastfm == "" {
		panic("No LastFM api key provided")
	}

	ctx := context.Background()

	s := sse.NewSSE(ctx)
	l := lastfm.NewLastFM(ctx, f.lastfm)
	h := handler.NewHandler(s, l, f.dataPath, f.secret)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestSize(maxBodySize))

	routes.RegisterAssetsRoutes(h, r)
	routes.RegisterPagesRoutes(h, r)
	routes.RegisterBlogRoutes(h, r)
	routes.RegisterChatRoutes(h, r)

	listen := fmt.Sprintf("%s:%d", f.address, f.port)

	slog.Info("Listening on: " + listen)
	http.ListenAndServe(listen, r)
}

func parseFlags() *flags {
	var f flags

	flag.StringVar(&f.secret, "secret", "", "Secret to use for various secret tasks")
	flag.StringVar(&f.lastfm, "last-fm", "", "Key to use for LastFM")
	flag.StringVar(&f.address, "address", "0.0.0.0", "Address to listen on")
	flag.IntVar(&f.port, "port", 3000, "Port to listen on")
	flag.StringVar(&f.dataPath, "data", "data", "Path to load app data from")

	flag.Parse()

	return &f
}
