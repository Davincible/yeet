package entrypoint

import (
	"crypto/tls"
	"time"
)

// Default config options.
const (
	DefaultAllowInsecure        = false
	DefaultMaxConcurrentStreams = 250
	DefaultReadTimeout          = 5 * time.Second
	DefaultWriteTimeout         = 5 * time.Second
	DefaultIdleimeout           = 5 * time.Second
)

// Option is a functional option to provide custom values to the config.
type Option func(*Config)

// Config provides options to the entrypoint.
type Config struct {
	Address string

	CertFile string
	KeyFile  string

	// Insecure will create an HTTP server without TLS, for insecure connections.
	Insecure bool
	// MaxConcurrentStreams for HTTP2.
	MaxConcurrentStreams int
	// TLS config, if none is provided self-signed certificates will be generated.
	TLS *tls.Config
	// AllowH2C allows h2c connections; HTTP2 without TLS.
	AllowH2C bool
	// HTTP3 dicates whether to also accept HTTP3.
	HTTP3 bool

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewEntrypointConfig will create a new config with default values for the entrypoint.
func NewEntrypointConfig(options ...Option) Config {
	cfg := Config{
		Insecure:             false,
		MaxConcurrentStreams: DefaultMaxConcurrentStreams,
		AllowH2C:             false,
		HTTP3:                false,
		ReadTimeout:          DefaultReadTimeout,
		WriteTimeout:         DefaultWriteTimeout,
		IdleTimeout:          DefaultIdleimeout,
	}

	cfg.ApplyOptions(options...)

	return cfg
}

func (c *Config) ApplyOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}

func WithAddress(address string) Option {
	return func(c *Config) {
		c.Address = address
	}
}

func WithTLSFile(certfile, keyfile string) Option {
	return func(c *Config) {
		c.CertFile = certfile
		c.KeyFile = keyfile
	}
}

func WithTLS(tlsConfig *tls.Config) Option {
	return func(c *Config) {
		c.TLS = tlsConfig
	}
}

func WithInsecure() Option {
	return func(c *Config) {
		c.Insecure = true
	}
}

// TODO: other options and option comments
