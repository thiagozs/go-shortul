package api

import (
	"log/slog"

	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/handler"
	"github.com/thiagozs/go-shorturl/infra/database"
	"github.com/thiagozs/go-shorturl/middleware"
)

type Options func(*APIParams) error

type APIParams struct {
	db         *database.Database
	middleware *middleware.Middleware
	handlers   *handler.Handler
	logger     *slog.Logger
	config     *config.Config
	domain     string
	port       string
	host       string
	https      bool
}

func newApiParams(opts ...Options) (*APIParams, error) {
	params := &APIParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func WithDB(db *database.Database) Options {
	return func(p *APIParams) error {
		p.db = db
		return nil
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(p *APIParams) error {
		p.logger = logger
		return nil
	}
}

func WithPort(port string) Options {
	return func(p *APIParams) error {
		p.port = port
		return nil
	}
}

func WithDomain(domain string) Options {
	return func(p *APIParams) error {
		p.domain = domain
		return nil
	}
}

func WithHost(host string) Options {
	return func(p *APIParams) error {
		p.host = host
		return nil
	}
}

func WithHTTPS(https bool) Options {
	return func(p *APIParams) error {
		p.https = https
		return nil
	}
}

func WithMiddleware(middleware *middleware.Middleware) Options {
	return func(p *APIParams) error {
		p.middleware = middleware
		return nil
	}
}

func WithHandlers(handlers *handler.Handler) Options {
	return func(p *APIParams) error {
		p.handlers = handlers
		return nil
	}
}

func WithConfig(config *config.Config) Options {
	return func(p *APIParams) error {
		p.config = config
		return nil
	}
}

// getters -----

func (p *APIParams) DB() *database.Database {
	return p.db
}

func (p *APIParams) Logger() *slog.Logger {
	return p.logger
}

func (p *APIParams) Domain() string {
	return p.domain
}

func (p *APIParams) Port() string {
	return p.port
}

func (p *APIParams) Host() string {
	return p.host
}

func (p *APIParams) HTTPS() bool {
	return p.https
}

func (p *APIParams) Middleware() *middleware.Middleware {
	return p.middleware
}

func (p *APIParams) Handlers() *handler.Handler {
	return p.handlers
}

func (p *APIParams) Config() *config.Config {
	return p.config
}

// setters -----

func (p *APIParams) SetDB(db *database.Database) {
	p.db = db
}

func (p *APIParams) SetLogger(logger *slog.Logger) {
	p.logger = logger
}

func (p *APIParams) SetDomain(domain string) {
	p.domain = domain
}

func (p *APIParams) SetPort(port string) {
	p.port = port
}

func (p *APIParams) SetHost(host string) {
	p.host = host
}

func (p *APIParams) SetHTTPS(https bool) {
	p.https = https
}

func (p *APIParams) SetMiddleware(middleware *middleware.Middleware) {
	p.middleware = middleware
}

func (p *APIParams) SetHandlers(handlers *handler.Handler) {
	p.handlers = handlers
}

func (p *APIParams) SetConfig(config *config.Config) {
	p.config = config
}
