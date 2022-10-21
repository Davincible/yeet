package http

import (
	"bytes"
	"io"
	"net/http"

	"http-poc/http/headers"
	"http-poc/http/utils/header"
)

func (s *Server) decodeBody(w http.ResponseWriter, r *http.Request, in any) (string, error) {
	ctHeader := r.Header.Get(headers.ContentType)
	contentType, err := header.GetContentType(ctHeader)
	if err != nil {
		return "", err
	}

	aHeader := r.Header.Get(headers.AcceptEncoding)
	accept := header.GetAcceptType(s.codecs, aHeader, contentType)
	w.Header().Set(headers.ContentType, accept)

	codec, ok := s.codecs[contentType]
	if !ok {
		return "", ErrContentTypeNotSupported
	}

	var b io.Reader

	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query().Encode()
		b = bytes.NewBufferString(query)

		contentType = headers.FormContentType
	default:
		b = r.Body
	}

	if err := codec.NewDecoder(b).Decode(in); err != nil {
		return "", err
	}

	return accept, nil
}

func (s *Server) encodeBody(w http.ResponseWriter, v any) error {
	contentType := w.Header().Get(headers.ContentType)
	if len(contentType) > 0 {
		contentType = headers.JSONContentType
	}

	codec, ok := s.codecs[contentType]
	if !ok {
		return ErrContentTypeNotSupported
	}

	if err := codec.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}
