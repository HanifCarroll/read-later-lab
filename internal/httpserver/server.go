package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go-embedded-js/internal/app"
)

type itemService interface {
	ListItems(context.Context) ([]app.Item, error)
	CreateItem(context.Context, app.CreateItemInput) (app.Item, error)
	UpdateStatus(context.Context, string, app.ItemStatus) (app.Item, error)
}

type Config struct {
	Service        itemService
	Assets         fs.FS
	DevSvelteProxy string
}

func New(config Config) http.Handler {
	s := &server{service: config.Service, assets: config.Assets, devSvelteProxy: config.DevSvelteProxy}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.home)
	mux.HandleFunc("GET /api/items", s.apiList)
	mux.HandleFunc("POST /api/items", s.apiCreate)
	mux.HandleFunc("POST /api/items/{id}/status", s.apiStatus)
	mux.Handle("GET /app", s.svelteApp())
	mux.Handle("GET /app/", s.svelteApp())
	return mux
}

type server struct {
	service        itemService
	assets         fs.FS
	devSvelteProxy string
}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app", http.StatusFound)
}

func (s *server) apiList(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListItems(r.Context())
	respond(w, items, err)
}

func (s *server) apiCreate(w http.ResponseWriter, r *http.Request) {
	var input app.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.service.CreateItem(r.Context(), input)
	respond(w, item, err)
}

func (s *server) apiStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status app.ItemStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.service.UpdateStatus(r.Context(), r.PathValue("id"), input.Status)
	respond(w, item, err)
}

func (s *server) svelteApp() http.Handler {
	if s.devSvelteProxy != "" {
		return s.svelteDevProxy()
	}

	fileServer := http.FileServer(http.FS(s.assets))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/app/")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(s.assets, path); err != nil {
			path = "index.html"
		}
		if path == "index.html" {
			content, err := fs.ReadFile(s.assets, path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(content)
			return
		}

		request := r.Clone(r.Context())
		request.URL.Path = "/" + path
		fileServer.ServeHTTP(w, request)
	})
}

func (s *server) svelteDevProxy() http.Handler {
	target, err := url.Parse(s.devSvelteProxy)
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "invalid DEV_SVELTE_PROXY: "+err.Error(), http.StatusInternalServerError)
		})
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func respond(w http.ResponseWriter, value any, err error) {
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func respondError(w http.ResponseWriter, status int, err error) {
	if errors.Is(err, context.Canceled) {
		status = http.StatusRequestTimeout
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
