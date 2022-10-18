package main

import (
	"context"
	"fmt"
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

func NewHandler[Tin any, Tout any](f func(context.Context, *Tin) (*Tout, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := new(Tin)

		out, _ := f(r.Context(), in)

		w.Write([]byte("hi"))
	}
}

func main() {
	handler := NewHandler[RequestA, ResponseA](Echo)
	fmt.Println(handler)
}
