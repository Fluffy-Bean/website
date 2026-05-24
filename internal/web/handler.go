package web

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"git.leggy.dev/Fluffy/Website/internal/broker"
	"git.leggy.dev/Fluffy/Website/internal/jwt"
	"git.leggy.dev/Fluffy/Website/internal/lastfm"
	"git.leggy.dev/Fluffy/Website/internal/sse"
)

const (
	cookieName     = "I-Heart-Maned-Wolves"
	cookieLifespan = 12 * time.Hour
)

type Data map[string]interface{}

type Handler struct {
	SSE       *sse.SSE
	LastFM    *lastfm.LastFM
	Events    *broker.Broker
	createdAt time.Time
	dataPath  string
	secret    string
}

func NewHandler(s *sse.SSE, l *lastfm.LastFM, dataPath, secret string) *Handler {
	return &Handler{
		SSE:       s,
		LastFM:    l,
		Events:    broker.NewBroker(),
		createdAt: time.Now().UTC(),
		dataPath:  dataPath,
		secret:    secret,
	}
}

func (h *Handler) DataPath(dir string) string {
	return path.Join(h.dataPath, dir)
}

func (h *Handler) parseTemplatesDir(dir string, templ *template.Template) (*template.Template, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("read templates dir", "error", err)

		return nil, err
	}

	for _, file := range files {
		templ, err = templ.ParseFiles(path.Join(dir, file.Name()))
		if err != nil {
			slog.Error("parse template file", "error", err)

			return nil, err
		}
	}

	return templ, nil
}

func (h *Handler) Template(w http.ResponseWriter, r *http.Request, path string, vars Data) {
	var err error

	templ := template.New("layout.html").Funcs(templateFuncs)

	templ, err = templ.ParseFiles("templates/layout.html", path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	templ, err = h.parseTemplatesDir("templates/blocks", templ)
	if err != nil {
		slog.Error("get block templates", "error", err)

		return
	}

	templ, err = h.parseTemplatesDir("templates/partials", templ)
	if err != nil {
		slog.Error("get partial templates", "error", err)

		return
	}

	templ, err = h.parseTemplatesDir("templates/blocks", templ)
	if err != nil {
		slog.Error("get block templates", "error", err)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = templ.Execute(w, map[string]interface{}{
		"URL":  tmplURLVars(r),
		"Vars": vars,
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

	if claims.Iat < h.createdAt.UTC().Unix() {
		return nil, fmt.Errorf("token is too old")
	}

	return claims, nil
}
