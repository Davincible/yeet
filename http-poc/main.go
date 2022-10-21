package main

import (
	"context"
	"net/http"
)

// func (s *httpserver) NewHandler[T any](f func(context.Context, T) (any, error)) {
//
// }

type RequestA struct {
	msg string
}

type ResponseA struct {
	msg string
}

func Echo(ctx context.Context, in *RequestA) (*ResponseA, error) {
	return &ResponseA{msg: "Hello " + in.msg}, nil
}

func main() {
	// handler := NewHandler[RequestA, ResponseA](Echo)
	// fmt.Println(handler)
}

func WriteError(rsp http.ResponseWriter, err error) {
	// TODO: implement this properly
	rsp.WriteHeader(500)
	_, _ = rsp.Write([]byte(err.Error()))
}
