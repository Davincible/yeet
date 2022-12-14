// Package entrypoint provides the entrypoint type used to attach servers
// to an address.
package entrypoint

import (
	"context"
	"crypto/tls"
	"net"

	"http-poc/http/router/router"
	mip "http-poc/http/utils/ip"
	mtcp "http-poc/http/utils/tcp"
	mtls "http-poc/http/utils/tls"
	mudp "http-poc/http/utils/udp"

	"github.com/pkg/errors"

	"http-poc/logger"
)

// Entrypoint represents a listener on one address. You can create multiple
// entrypoints for multiple addresses and ports. This is e.g. useful if you
// want to listen on multiple interfaces, or multiple ports in parallel, even
// with the same handler.
type Entrypoint struct {
	Config Config
	logger logger.Logger

	listenerUDP net.PacketConn
	listenerTCP net.Listener

	httpServer  *httpServer
	http3Server *http3server
}

// NewEntrypoint creates a new entrypoint for a single address. You can create
// multiple entrypoints for multiple addresses and ports. One entrypoint
// can serve a HTTP1, HTTP2 and HTTP3 server. If you enable HTTP3 it will listen
// on both TCP and UDP on the same port.
func NewEntrypoint(router router.Router, logger logger.Logger, config Config, options ...Option) (*Entrypoint, error) {
	config.ApplyOptions(options...)

	if err := mip.ValidateAddress(config.Address); err != nil {
		return nil, err
	}

	entrypoint := Entrypoint{
		Config: config,
		logger: logger,
	}

	var err error

	entrypoint.Config.TLS, err = entrypoint.setupTLS()
	if err != nil {
		return nil, err
	}

	entrypoint.httpServer, err = entrypoint.newHTTPServer(router)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP server")
	}

	if !entrypoint.Config.HTTP3 {
		return &entrypoint, nil
	}

	entrypoint.http3Server, err = entrypoint.newHTTP3Server()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP3 server")
	}

	return &entrypoint, nil
}

// Start will create the listeners and start the server on the entrypoint.
func (e *Entrypoint) Start() error {
	var err error

	// e.logger.Debug("Starting all HTTP entrypoints")

	e.listenerTCP, err = mtcp.BuildListenerTCP(e.Config.Address, e.Config.TLS)
	if err != nil {
		return err
	}

	go func() {
		if err = e.httpServer.Start(e.listenerTCP); err != nil {
			e.logger.Errorf("Failed to start HTTP server: %v", err)
		}
	}()

	if !e.Config.HTTP3 {
		return nil
	}

	// Listen on the same UDP port as TCP for HTTP3
	e.listenerUDP, err = mudp.BuildListenerUDP(e.Config.Address)
	if err != nil {
		return errors.Wrap(err, "failed to start UDP listener")
	}

	go func() {
		if err := e.http3Server.Start(e.listenerUDP); err != nil {
			e.logger.Errorf("Failed to start HTTP3 server: %v", err)
		}
	}()

	return nil
}

type stopper interface {
	Stop(context.Context) error
}

// Stop will stop all servers and close the listeners.
func (e *Entrypoint) Stop(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)

	// e.logger.Debug("Stopping all HTTP entrypoints")

	c := 1
	if e.Config.HTTP3 {
		c++

		go func() {
			errChan <- e.http3Server.Stop(ctx)

			//nolint:errcheck
			_ = e.listenerUDP.Close()
		}()
	}

	go func(srv stopper, l net.Listener) {
		errChan <- srv.Stop(ctx)

		// TCP listener probably already closed, just as a double check.
		//nolint:errcheck
		_ = l.Close()
	}(e.httpServer, e.listenerTCP)

	var err error

	for i := 0; i < c; i++ {
		if nerr := <-errChan; nerr != nil {
			err = nerr
		}
	}

	return err
}

func (e *Entrypoint) setupTLS() (*tls.Config, error) {
	var (
		config *tls.Config
		err    error
	)

	// Load TLS from file
	if len(e.Config.CertFile) > 0 && len(e.Config.KeyFile) > 0 && e.Config.TLS == nil {
		config, err = mtls.LoadTLSConfig(e.Config.CertFile, e.Config.KeyFile)
		if err != nil {
			return config, errors.Wrap(err, "failed to load TLS config from files")
		}
	}

	// Generate self signed cert
	if !e.Config.Insecure && e.Config.TLS == nil {
		config, err = mtls.GenTlSConfig(e.Config.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate self signed certificate")
		}
	}

	return config, nil
}
