package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/thiagozs/go-shorturl/config"
)

type Handlers func(w http.ResponseWriter, r *http.Request)

type Middlewares func(http.Handler) http.Handler

type MiddlewaresFunc func(http.HandlerFunc) http.HandlerFunc

// Middlewares holds the logger and authentication token
type Middleware struct {
	params    *MiddlewareParams
	logger    *slog.Logger
	authToken string
}

func NewMiddleware(opts ...Options) (*Middleware, error) {
	params, err := newMiddlewareParams(opts...)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		params:    params,
		logger:    params.Logger(),
		authToken: params.Token(),
	}, nil
}

// loggingMiddleware is a middleware function that logs details about each request
func (m *Middleware) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.logger.Info("Received request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("referer", r.Referer()),
		)
		next.ServeHTTP(w, r)
		m.logger.Info("Processed request", slog.Duration("duration", time.Since(start)))
	}
}

// tokenAuthMiddleware is a middleware function that checks for a valid token in the request header
func (m *Middleware) TokenAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != m.authToken {
			m.logger.Warn("Unauthorized access attempt", slog.String("remote_addr", r.RemoteAddr))
			http.Error(w, "Forbidden: Invalid or missing token", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Let the request pass
		next.ServeHTTP(w, r)
	}
}

func (m Middleware) SugarM(middlewares []Middlewares, handler Handlers) http.Handler {
	var chain http.Handler
	for i, middleware := range middlewares {
		if i == 0 {
			chain = middleware(http.HandlerFunc(handler))
		} else {
			chain = middleware(chain)
		}
	}

	return chain
}

func (m Middleware) SugarMFunc(middlewares []MiddlewaresFunc, handler Handlers) http.HandlerFunc {
	var chain http.HandlerFunc
	for i, middleware := range middlewares {
		if i == 0 {
			chain = middleware(http.HandlerFunc(handler))
		} else {
			chain = middleware(chain)
		}
	}

	return chain
}

func (h *Middleware) SetConfig(cfg *config.Config) {
	h.authToken = cfg.GetToken()
	h.params.SetConfig(cfg)
}
