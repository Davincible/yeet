package http

import (
	"context"
	"fmt"
	"net/http"
)

// Convert proto function to http handler
func NewHandler[Tin any, Tout any](s ServerHTTP, f func(context.Context, *Tin) (*Tout, error)) http.HandlerFunc {
	srv := s.(*Server)

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: when to close the body?
		in := new(Tin)

		srv.decodeBody(w, r, in)

		out, err := f(r.Context(), in)
		if err != nil {
			WriteError(w, err)
			return
		}

		if err := srv.encodeBody(w, out); err != nil {
			WriteError(w, err)
			return
		}
	}
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, err.Error())
}
