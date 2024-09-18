package initialize

import (
	"log/slog"

	"github.com/thiagozs/go-shorturl/api"
	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/handler"
	"github.com/thiagozs/go-shorturl/infra/database"
	"github.com/thiagozs/go-shorturl/middleware"
)

type Options func(*InitializeParams) error

type InitializeParams struct {
	db         *database.Database
	handler    *handler.Handler
	middleware *middleware.Middleware
	logger     *slog.Logger
	config     *config.Config
	api        *api.API
}

func newInitializeParams(opts ...Options) (*InitializeParams, error) {
	params := &InitializeParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func WithDB(db *database.Database) Options {
	return func(p *InitializeParams) error {
		p.db = db
		return nil
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(p *InitializeParams) error {
		p.logger = logger
		return nil
	}
}

func WithHandler(handler *handler.Handler) Options {
	return func(p *InitializeParams) error {
		p.handler = handler
		return nil
	}
}

func WithConfig(config *config.Config) Options {
	return func(p *InitializeParams) error {
		p.config = config
		return nil
	}
}

func WithMiddleware(middleware *middleware.Middleware) Options {
	return func(p *InitializeParams) error {
		p.middleware = middleware
		return nil
	}
}

func WithAPI(api *api.API) Options {
	return func(p *InitializeParams) error {
		p.api = api
		return nil
	}
}

// getters -----

func (p *InitializeParams) GetDB() *database.Database {
	return p.db
}

func (p *InitializeParams) GetLogger() *slog.Logger {
	return p.logger
}

func (p *InitializeParams) GetHandlers() *handler.Handler {
	return p.handler
}

func (p *InitializeParams) GetConfig() *config.Config {
	return p.config
}

func (p *InitializeParams) GetMiddleware() *middleware.Middleware {
	return p.middleware
}

func (p *InitializeParams) GetAPI() *api.API {
	return p.api
}

// setters -----

func (p *InitializeParams) SetDB(db *database.Database) {
	p.db = db
}

func (p *InitializeParams) SetLogger(logger *slog.Logger) {
	p.logger = logger
}

func (p *InitializeParams) SetHandler(handler *handler.Handler) {
	p.handler = handler
}

func (p *InitializeParams) SetConfig(config *config.Config) {
	p.config = config
}

func (p *InitializeParams) SetMiddleware(middleware *middleware.Middleware) {
	p.middleware = middleware
}

func (p *InitializeParams) SetAPI(api *api.API) {
	p.api = api
}
