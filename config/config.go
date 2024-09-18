package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Host   string `env:"HOST"`
	Port   string `env:"PORT"`
	Domain string `env:"DOMAIN"`
	HTTPS  bool   `env:"HTTPS"`
	Local  bool   `env:"LOCAL"`
	Token  string `env:"SUPERSCRT"`
}

func NewConfig() (*Config, error) {

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *Config) GetHost() string {
	return c.Host
}

func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) GetDomain() string {
	return c.Domain
}

func (c *Config) GetHTTPS() bool {
	return c.HTTPS
}

func (c *Config) GetLocal() bool {
	return c.Local
}

func (c *Config) GetToken() string {
	return c.Token
}

// setters -----

func (c *Config) SetHost(host string) {
	c.Host = host
}

func (c *Config) SetPort(port string) {
	c.Port = port
}

func (c *Config) SetDomain(domain string) {
	c.Domain = domain
}

func (c *Config) SetHTTPS(https bool) {
	c.HTTPS = https
}

func (c *Config) SetLocal(local bool) {
	c.Local = local
}

func (c *Config) SetToken(token string) {
	c.Token = token
}
