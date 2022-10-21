package http

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"net/http"

	"go-micro.dev/v4/logger"

	"http-poc/handler"
	"http-poc/http/codec"
	"http-poc/http/router/chi"
	"http-poc/logger/zerolog"
)

func TestServerSimple(t *testing.T) {
	router := chi.ProvideChiRouter()
	logger := zerolog.ProvideZerologLogger(logger.WithOutput(os.Stdout), zerolog.WithDevelopmentMode())
	codecs := codec.ProvideCodecs(codec.ProvideDefaultCodecs()...)

	server, err := ProvideServerHTTP(router, codecs, logger, WithInsecure())
	if err != nil {
		t.Fatalf("failed to provide http server: %v", err)
	}

	// TODO: Normal request on TLS server, strange error
	h := new(handler.EchoHandler)
	router.Get("/echo", NewHandler(server, h.Call))
	router.Post("/echo", NewHandler(server, h.Call))

	if err := server.Start(); err != nil {
		t.Fatal("failed to start", err)
	}

	resp, err := http.Post("http://localhost:42069/echo", "application/json", bytes.NewBufferString(`{"name": "hi there!"}`))
	if err != nil {
		t.Fatal("failed to make post request", err)
	}
	logger.Infof("Response: %+v", resp)

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))

	if err := server.Stop(ctx); err != nil {
		t.Fatal("failed to stop", err)
	}
}
