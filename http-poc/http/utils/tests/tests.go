package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"

	"http-poc/http/headers"
	pb "http-poc/proto"
)

type ReqType int

// Request types.
const (
	TypeInsecure ReqType = iota + 1
	TypeHTTP2
	TypeHTTP3
	TypeH2C
)

// TestGetRequest makes a GET request to the echo endpoint.
func TestGetRequest(t testing.TB, addr string, reqT ReqType) error {
	name := "Alex"

	url := fmt.Sprintf("%s/echo?name=%s", addr, name)

	var (
		body []byte
		err  error
	)

	switch reqT {
	case TypeInsecure:
		body, err = makeGetReq(t, url, &http.Client{})
	case TypeHTTP2:
		body, err = makeSecureGetReq(t, url)
	case TypeHTTP3:
		body, err = makeHTTP3GetReq(t, url)
	case TypeH2C:
		body, err = makeH2CGetReq(t, url)
	}
	if err != nil {
		return err
	}

	if err := checkJSONResponse(body, name); err != nil {
		return err
	}

	return nil
}

// TestPostRequestJSON makes a POST request to the echo endpoint.
func TestPostRequestJSON(t testing.TB, addr string, reqT ReqType) error {
	name := "Alex"

	msg, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		t.Fatal("failed to marshall json", err)
	}

	addr += "/echo"
	ct := headers.JSONContentType

	var body []byte

	switch reqT {
	case TypeInsecure:
		body, err = makePostReq(t, addr, ct, msg, &http.Client{})
	case TypeHTTP2:
		body, err = makeSecurePostReq(t, addr, ct, msg)
	case TypeHTTP3:
		body, err = makeHTTP3PostReq(t, addr, ct, msg)
	case TypeH2C:
		body, err = makeH2CPostReq(t, addr, ct, msg)
	}
	if err != nil {
		return err
	}

	if err := checkJSONResponse(body, name); err != nil {
		return err
	}

	return nil
}

// TestPostRequestProto makes a POST request to the echo endpoint.
func TestPostRequestProto(t testing.TB, addr, ct string, reqT ReqType) error {
	name := "Alex"

	msg, err := proto.Marshal(&pb.CallRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}

	addr += "/echo"

	var body []byte

	switch reqT {
	case TypeInsecure:
		body, err = makePostReq(t, addr, ct, msg, &http.Client{})
	case TypeHTTP2:
		body, err = makeSecurePostReq(t, addr, ct, msg)
	case TypeHTTP3:
		body, err = makeHTTP3PostReq(t, addr, ct, msg)
	case TypeH2C:
		body, err = makeH2CPostReq(t, addr, ct, msg)
	}
	if err != nil {
		return err
	}

	if err := checkProtoResponse(body, name); err != nil {
		return err
	}
	return nil
}

// temporary test stuff
func TestTLSProto(t testing.TB, addr string) error {
	t.Log("Testing TLS")

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"HTTP/3.0", "My custom proto", "ClientGarbage"},
	})
	if err != nil {
		return fmt.Errorf("failed to dial TLS tcp connection: %w", err)
	}

	state := conn.ConnectionState()
	t.Log(state.NegotiatedProtocol)

	return nil
}

func checkJSONResponse(body []byte, name string) error {
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return errors.Wrap(err, "Failed to unmarhsal data")
	}

	if data["msg"] != "Hello "+name {
		return fmt.Errorf("request failed; expected different response than: %v", data)
	}

	return nil
}

func checkProtoResponse(body []byte, name string) error {
	var data pb.CallResponse
	if err := proto.Unmarshal(body, &data); err != nil {
		return errors.Wrap(err, "Failed to unmarhsal data")
	}

	if data.Msg != "Hello "+name {
		return fmt.Errorf("request failed; expected different response than: %v", data.Msg)
	}

	return nil
}

func makeH2CGetReq(t testing.TB, addr string) ([]byte, error) {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	return makeGetReq(t, addr, &client)
}

func makeH2CPostReq(t testing.TB, addr, ct string, msg []byte) ([]byte, error) {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	return makePostReq(t, addr, ct, msg, &client)
}

func makeSecureGetReq(t testing.TB, addr string) ([]byte, error) {
	client := http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	return makeGetReq(t, addr, &client)
}

func makeSecurePostReq(t testing.TB, addr, ct string, msg []byte) ([]byte, error) {
	client := http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	return makePostReq(t, addr, ct, msg, &client)
}

func makeHTTP3GetReq(t testing.TB, addr string) ([]byte, error) {
	client := http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	return makeGetReq(t, addr, &client)
}

func makeHTTP3PostReq(t testing.TB, addr, ct string, msg []byte) ([]byte, error) {
	client := http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	return makePostReq(t, addr, ct, msg, &client)
}

func makeGetReq(t testing.TB, addr string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// logResponse(t, resp, body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}

	return body, nil
}

func makePostReq(t testing.TB, addr, ct string, data []byte, client *http.Client) ([]byte, error) {
	resp, err := client.Post(addr, ct, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// logResponse(t, resp, body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Post request failed")
	}

	return body, nil
}

func logResponse(tb testing.TB, resp *http.Response, body []byte) {
	tb.Logf(
		"[%+v] Status: %v, \n\tProto: %v, ConentType: %v, Length: %v, \n\tTransferEncoding: %v, Uncompressed: %v, \n\tBody: %v",
		resp.Request.Method, resp.Status, resp.Proto, resp.Header.Get("Content-Type"), resp.ContentLength, resp.TransferEncoding, resp.Uncompressed, string(body),
	)
}
