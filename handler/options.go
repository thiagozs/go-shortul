package handler

import (
	"log/slog"

	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/infra/database"
)

type Options func(*HandlerParams) error

type HandlerParams struct {
	store  *database.Database
	logger *slog.Logger
	config *config.Config
	domain string
	port   string
	host   string
	local  bool
	https  bool
}

func newHandlerParams(opts ...Options) (*HandlerParams, error) {
	params := &HandlerParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func WithStore(store *database.Database) Options {
	return func(p *HandlerParams) error {
		p.store = store
		return nil
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(p *HandlerParams) error {
		p.logger = logger
		return nil
	}
}

func WithPort(port string) Options {
	return func(p *HandlerParams) error {
		p.port = port
		return nil
	}
}

func WithHost(host string) Options {
	return func(p *HandlerParams) error {
		p.host = host
		return nil
	}
}

func WithLocal(local bool) Options {
	return func(p *HandlerParams) error {
		p.local = local
		return nil
	}
}

func WithHTTPS(https bool) Options {
	return func(p *HandlerParams) error {
		p.https = https
		return nil
	}
}

func WithDomain(domain string) Options {
	return func(p *HandlerParams) error {
		p.domain = domain
		return nil
	}
}

func WithConfig(config *config.Config) Options {
	return func(p *HandlerParams) error {
		p.config = config
		return nil
	}
}

// getters -----

func (p *HandlerParams) Store() *database.Database {
	return p.store
}

func (p *HandlerParams) Logger() *slog.Logger {
	return p.logger
}

func (p *HandlerParams) Port() string {
	return p.port
}

func (p *HandlerParams) Host() string {
	return p.host
}

func (p *HandlerParams) Local() bool {
	return p.local
}

func (p *HandlerParams) HTTPS() bool {
	return p.https
}

func (p *HandlerParams) Domain() string {
	return p.domain
}

func (p *HandlerParams) Config() *config.Config {
	return p.config
}

// setters -----

func (p *HandlerParams) SetStore(store *database.Database) {
	p.store = store
}

func (p *HandlerParams) SetLogger(logger *slog.Logger) {
	p.logger = logger
}

func (p *HandlerParams) SetPort(port string) {
	p.port = port
}

func (p *HandlerParams) SetHost(host string) {
	p.host = host
}

func (p *HandlerParams) SetLocal(local bool) {
	p.local = local
}

func (p *HandlerParams) SetHTTPS(https bool) {
	p.https = https
}

func (p *HandlerParams) SetDomain(domain string) {
	p.domain = domain
}

func (p *HandlerParams) SetConfig(config *config.Config) {
	p.config = config
}
