package web

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"git.leggy.dev/Fluffy/Website/data"
	"git.leggy.dev/Fluffy/Website/internal/broker"
	"git.leggy.dev/Fluffy/Website/internal/jwt"
	"git.leggy.dev/Fluffy/Website/internal/lastfm"
	"git.leggy.dev/Fluffy/Website/internal/sse"
	"git.leggy.dev/Fluffy/Website/templates"
)

const (
	cookieName     = "I-Heart-Maned-Wolves"
	cookieLifespan = 12 * time.Hour
)

type Data map[string]interface{}

type Handler struct {
	secret    string
	SSE       *sse.SSE
	LastFM    *lastfm.LastFM
	Events    *broker.Broker
	StartTime time.Time
	BuildTime time.Time
}

func NewHandler(s *sse.SSE, l *lastfm.LastFM, secret string) *Handler {
	return &Handler{
		secret:    secret,
		SSE:       s,
		LastFM:    l,
		Events:    broker.NewBroker(),
		StartTime: time.Now().UTC(),
		BuildTime: time.Now().UTC(),
	}
}

func (h *Handler) ReadDataFile(file string) ([]byte, error) {
	if _, err := os.Stat("data"); err != nil {
		return data.Dir.ReadFile(file)
	}
	return os.ReadFile(path.Join("data", file))
}

func (h *Handler) ReadDataDir(dir string) ([]os.DirEntry, error) {
	if _, err := os.Stat("data"); err != nil {
		return data.Dir.ReadDir(dir)
	}
	return os.ReadDir(path.Join("data", dir))
}

func (h *Handler) Uptime() time.Duration {
	return time.Now().Sub(h.StartTime)
}

func (h *Handler) templateParse(templ *template.Template, patterns ...string) (*template.Template, error) {
	if _, err := os.Stat("templates"); err != nil {
		return templ.ParseFS(templates.Dir, patterns...)
	}

	var err error
	for _, pattern := range patterns {
		templ, err = templ.ParseFiles(path.Join("templates", pattern))
		if err != nil {
			return nil, err
		}
	}

	return templ, nil
}

func (h *Handler) ReadTemplatesDir(dir string) ([]os.DirEntry, error) {
	if _, err := os.Stat("templates"); err != nil {
		return templates.Dir.ReadDir(dir)
	}
	return os.ReadDir(path.Join("templates", dir))
}

func (h *Handler) Template(w http.ResponseWriter, r *http.Request, page string, vars Data) {
	var templ *template.Template
	var err error

	templ = template.New("layout.html").Funcs(templateFuncs)

	templ, err = h.templateParse(templ, "layout.html", path.Join("pages", page))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	blocks, err := h.ReadTemplatesDir("blocks")
	if err != nil {
		slog.Error("read templates dir", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	for _, block := range blocks {
		templ, err = h.templateParse(templ, path.Join("blocks", block.Name()))
		if err != nil {
			slog.Error("parse template file", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = templ.Execute(w, map[string]interface{}{
		"URL":       tmplURLVars(r),
		"Time":      tmplTimeVars(),
		"StartTime": h.StartTime,
		"BuildTime": h.BuildTime,
		"Uptime":    h.Uptime(),
		"Vars":      vars,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) SetToken(w http.ResponseWriter, id int64, name string) error {
	now := time.Now().UTC()

	claims := jwt.Claims{
		Sub:  id,
		Name: name,
		Iat:  now.Unix(),
	}

	token, err := jwt.New(h.secret).Encode(claims)
	if err != nil {
		return fmt.Errorf("encode token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Expires:  now.Add(cookieLifespan),
	})

	return nil
}

func (h *Handler) GetToken(r *http.Request) (*jwt.Claims, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == cookieName {
			cookie = c

			break
		}
	}

	if cookie == nil {
		return nil, fmt.Errorf("no cookie set")
	}

	claims, err := jwt.New(h.secret).VerifyAndDecode(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("decoding jwt token: %w", err)
	}

	if claims.Iat < h.StartTime.UTC().Unix() {
		return nil, fmt.Errorf("token is too old")
	}

	return claims, nil
}
