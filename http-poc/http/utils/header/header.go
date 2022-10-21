// Package header implements header manipulation utilities
package header

import (
	"fmt"
	"mime"
	"strings"

	"http-poc/http/codec"
)

// GetContentType parses the content type from the header value.
func GetContentType(header string) (string, error) {
	ct, _, err := mime.ParseMediaType(header)
	if err != nil {
		// TODO: return custom error, log this error
		return "", err
	}

	return ct, nil
}

// GetAcceptType parses the Accept header and checks against the available codecs
// to find a matching content type.
func GetAcceptType(c codec.Codecs, acceptHeader string, contentType string) string {
	accept := contentType

	acceptSlice := strings.Split(acceptHeader, ",")
	for _, acceptType := range acceptSlice {
		ct, _, err := mime.ParseMediaType(acceptType)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Check if we have a codec for the content type
		if _, ok := c[ct]; ok {
			accept = ct
		}
	}

	return accept
}