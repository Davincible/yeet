package http

import (
	"fmt"

	"http-poc/http/entrypoint"
)

const (
	DefaultAddress = "0.0.0.0:42069"
)

// Option is a functional HTTP server option.
type Option func(*Config)

// Entrypoint is a list of entrypoint options used to facilitate the creation
// of one entrypoint, with all options provided. You MUST atleast add the address
// to listen on WithAddress, and can optionally specify more options to use.
// For all options not specified, default values will be used.
type Entrypoint []entrypoint.Option

// Config is the server config. It contains the list of addresses on which
// entrypoints will be created, and the default config used for each entrypoint.
type Config struct {
	EntrypointDefaults entrypoint.Config
	Entrypoints        []Entrypoint
}

// NewConfig creates a new server config with default values.
// To customise the options pas in a list of options.
func NewConfig(options ...Option) Config {
	cfg := Config{
		EntrypointDefaults: entrypoint.NewEntrypointConfig(),
		Entrypoints:        make([]Entrypoint, 0, 1),
	}

	cfg.ApplyOptions(options...)

	if len(cfg.Entrypoints) == 0 {
		fmt.Println("adding")
		e := Entrypoint{entrypoint.WithAddress(DefaultAddress)}
		cfg.Entrypoints = append(cfg.Entrypoints, e)
	}

	fmt.Println(cfg.Entrypoints, len(cfg.Entrypoints))

	return cfg
}

// ApplyOptions takes a list of options and applies them to the current config.
func (c *Config) ApplyOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}

// WithEntrypointOptions takes a list of entrypoint.Option to apply to the
// EntrypointDefaults, to use as defaults when new entrypoints are created.
// See entrypoint.Config for more details about the possible options.
func WithEntrypointOptions(options ...entrypoint.Option) Option {
	return func(c *Config) {
		c.EntrypointDefaults.ApplyOptions(options...)
	}
}

// WithEntrypointDefaults takes a complete entrypoint.Config to use as default
// when creating new entrypoints. Only use this if you want to specify (almost)
// all values in the entrypoint.Config. Otherwise use WithEntrypointOptions to
// specify a list of options.
func WithEntrypointDefaults(defaults entrypoint.Config) Option {
	return func(c *Config) {
		c.EntrypointDefaults = defaults
	}
}

// WithAddress takes a list of addresses to listen on.
// This is an alias for WithEntrypoints, specifying only the address.
//
// It will create an Entrypoint for each address, attaching an HTTP server on each.
// Each entrypoint will be created with the Config.DefaultEntryPointConfig.
// If you want to change the defaults used for the creation of each entrypoint
// specify them WithEntrypointOptions, if you want to create custom entrypoints
// with specific non-default options, use WithEntrypoints.
//
// To listen on all interfaces on one specific port use use the ":<port>" notation.
// To listen on a specific interface use the "<IP>:<port>" notation.
func WithAddress(addrs ...string) Option {
	return func(c *Config) {
		for _, addr := range addrs {
			e := Entrypoint{entrypoint.WithAddress(addr)}
			c.Entrypoints = append(c.Entrypoints, e)
		}
	}
}

// WithEntrypoints is used to provide a list of entrypoints to create, will
// custom options only applied to each entypoint individually.
// You MUST atleast provide an address to listen on WithAddress.
//
// Example:
//
//	WithEntrypoints([]http.Entrypoint{
//	 {entrypoint.WithAddress(":8080")},
//	 {entrypoint.WithAddress(":8081"), WithHTTP3()},
//	 {entrypoint.WithAddress(":8082"), WithInsecure(), WithWriteTimeout(...)},
//	})
func WithEntrypoints(entrypoints ...Entrypoint) Option {
	return func(c *Config) {
		c.Entrypoints = append(c.Entrypoints, entrypoints...)
	}
}

func WithInsecure() Option {
	return WithEntrypointOptions(entrypoint.WithInsecure())
}
