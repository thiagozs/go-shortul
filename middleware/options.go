package middleware

import (
	"log/slog"

	"github.com/thiagozs/go-shorturl/config"
)

type Options func(*MiddlewareParams) error

type MiddlewareParams struct {
	config *config.Config
	logger *slog.Logger
	token  string
}

func newMiddlewareParams(opts ...Options) (*MiddlewareParams, error) {
	params := &MiddlewareParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func WithLogger(logger *slog.Logger) Options {
	return func(p *MiddlewareParams) error {
		p.logger = logger
		return nil
	}
}

func WithToken(token string) Options {
	return func(p *MiddlewareParams) error {
		p.token = token
		return nil
	}
}

func WithConfig(config *config.Config) Options {
	return func(p *MiddlewareParams) error {
		p.config = config
		return nil
	}
}

func (p *MiddlewareParams) Config() *config.Config {
	return p.config
}

func (p *MiddlewareParams) Logger() *slog.Logger {
	return p.logger
}

func (p *MiddlewareParams) Token() string {
	return p.token
}

func (p *MiddlewareParams) SetLogger(logger *slog.Logger) {
	p.logger = logger
}

func (p *MiddlewareParams) SetToken(token string) {
	p.token = token
}

func (p *MiddlewareParams) SetConfig(config *config.Config) {
	p.config = config
}
