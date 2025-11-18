package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Server struct {
	router        *Router
	logger        *slog.Logger
	db            *sql.DB
	staticPrefix  string
	staticHandler http.Handler
}

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func New() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	router := &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}

	s := &Server{
		router: router,
		logger: logger,
	}

	return s
}

func NewWithDB(db *sql.DB) *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	router := &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}

	s := &Server{
		router: router,
		logger: logger,
		db:     db,
	}

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.staticHandler != nil && strings.HasPrefix(r.URL.Path, s.staticPrefix) {
		s.staticHandler.ServeHTTP(w, r)
		return
	}
	s.recoveryMiddleware(s.middleware(s.router)).ServeHTTP(w, r)
}

func (s *Server) GetRouter() *Router {
	return s.router
}

func (s *Server) ServeStatic(prefix, dir string) {
	if prefix == "" {
		return
	}

	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	s.staticPrefix = prefix
	s.staticHandler = http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
}

func (r *Router) HandleFunc(pattern, method string, handler http.HandlerFunc) {
	if r.routes[pattern] == nil {
		r.routes[pattern] = make(map[string]http.HandlerFunc)
	}
	r.routes[pattern][method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// First try exact match
	if routes, exists := r.routes[path]; exists {
		if handler, exists := routes[method]; exists {
			handler(w, req)
			return
		}
	}

	// Then try pattern matching for routes with parameters
	for pattern, routes := range r.routes {
		if r.matchesPattern(pattern, path) {
			if handler, exists := routes[method]; exists {
				handler(w, req)
				return
			}
		}
	}

	http.NotFound(w, req)
}

func (r *Router) matchesPattern(pattern, path string) bool {
	// Simple pattern matching for {id} parameters
	// Convert pattern like "/api/users/{id}/stats" to regex-like matching
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, patternPart := range patternParts {
		// Handle empty parts (leading/trailing slashes)
		if patternPart == "" && pathParts[i] == "" {
			continue
		}
		if patternPart == "" || pathParts[i] == "" {
			return false
		}
		// Check if this is a parameter like {id}
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			// This is a parameter, match any non-empty string
			if pathParts[i] == "" {
				return false
			}
			continue
		}
		// Exact match for non-parameter parts
		if patternPart != pathParts[i] {
			return false
		}
	}

	return true
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return s.corsMiddleware(s.loggingMiddleware(next))
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.logger.Error("panic recovered",
					"error", rec,
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		ctx := context.WithValue(r.Context(), "logger", s.logger)
		r = r.WithContext(ctx)

		s.logger.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		s.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
