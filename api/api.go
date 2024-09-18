package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/handler"
	"github.com/thiagozs/go-shorturl/middleware"
)

type API struct {
	params *APIParams
	server *http.Server
	http   *http.ServeMux
	logger *slog.Logger
	md     *middleware.Middleware
	hd     *handler.Handler
}

func NewApi(opts ...Options) (*API, error) {
	params, err := newApiParams(opts...)
	if err != nil {
		return nil, err
	}

	return &API{
		params: params,
		http:   http.NewServeMux(),
		logger: params.Logger(),
		md:     params.Middleware(),
		hd:     params.Handlers(),
	}, nil
}

func (a *API) RegisterServer() error {
	a.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", a.params.Host(), a.params.Port()),
		Handler: a.http,
	}
	return nil
}

func (a *API) RegisterEndPoints() error {
	// Define the middleware functions to use auth
	midAuth := []middleware.MiddlewaresFunc{
		a.md.CORS,
		a.md.Logging,
		a.md.TokenAuth,
	}

	// Define the middleware functions to use common
	midcommon := []middleware.MiddlewaresFunc{
		a.md.CORS,
		a.md.Logging,
	}

	endpoints := map[string]http.HandlerFunc{
		"/shorten": a.md.SugarMFunc(midAuth, a.hd.ShortenHandler),
		"/update":  a.md.SugarMFunc(midAuth, a.hd.UpdateHandler),
		"/flush":   a.md.SugarMFunc(midAuth, a.hd.FlushHandler),
		"/backup":  a.md.SugarMFunc(midAuth, a.hd.BackupHandler),
		"/import":  a.md.SugarMFunc(midAuth, a.hd.ImportHandler),
		"/stats":   a.md.SugarMFunc(midAuth, a.hd.StatsHandler),
		"/":        a.md.SugarMFunc(midcommon, a.hd.RedirectHandler),
		"/health":  a.md.SugarMFunc(midcommon, a.hd.HealthHandler),
	}

	for path, handler := range endpoints {
		a.http.HandleFunc(path, handler)
	}

	return nil
}

func (a *API) Start() {
	go func(*API) {
		a.logger.Info("Server started")
		if err := a.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			a.logger.Error("listen and serve", "error", err)
		}
		a.logger.Info("Server stopped")
	}(a)
}

func (a *API) Shutdown() error {
	return a.server.Shutdown(context.Background())
}

func (a *API) SetConfigByFlags(cfg *config.Config) {
	a.params.SetConfig(cfg)
	a.params.Handlers().SetConfig(cfg)
	a.params.Middleware().SetConfig(cfg)

	if a.params.Host() == "" {
		a.params.SetHost(a.params.Config().GetHost())
	}

	if a.params.Port() == "" {
		a.params.SetPort(a.params.Config().GetPort())
	}

}
