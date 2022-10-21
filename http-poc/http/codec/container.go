package codec

import (
	"errors"

	"http-poc/http/codec/form"

	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Codec interface {
	runtime.Marshaler
}

type JSONPb struct {
	runtime.JSONPb

	// contentType is used to overwrite the content type used by the encoder.
	// This is useful when one encoder encodes for multiple content types,
	// and you have to create a sperate instance for each one.
	contentType string
}

func (j *JSONPb) ContentType(_ any) string {
	return j.contentType
}

type Proto struct {
	runtime.ProtoMarshaller

	// contentType is used to overwrite the content type used by the encoder.
	// This is useful when one encoder encodes for multiple content types,
	// and you have to create a sperate instance for each one.
	contentType string
}

func (p *Proto) ContentType(_ any) string {
	return p.contentType
}

type Codecs map[string]Codec

type CodecRegistration struct {
	// One or more content types for which the codec is responsible
	ContentTypes []string
	Codec        Codec
}

var (
	CodecRegistry   = make(Codecs)
	ErrNotProtoType = errors.New("provided interface is not of type proto.Message")

	DefaultCodecSet = wire.NewSet(ProvideCodecJSON, ProvideCodecProto, ProvideCodecForm)
)

func ProvideCodecJSON() Codec {
	return &JSONPb{
		contentType: "application/json",
	}
}

func ProvideCodecProto() []Codec {
	return []Codec{
		&Proto{contentType: "application/octet-stream"},
		&Proto{contentType: "application/protobuf"},
		&Proto{contentType: "application/x-protobuf"},
	}
}

func ProvideCodecForm() Codec {
	return form.NewFormCodec()
}

func ProvideDefaultCodecs() []Codec {
	c := make([]Codec, 10)
	c = append(c, ProvideCodecForm())
	c = append(c, ProvideCodecJSON())
	c = append(c, ProvideCodecProto()...)

	return c
}

func ProvideCodecs(codecs ...Codec) map[string]Codec {
	m := make(map[string]Codec, len(codecs))

	for _, c := range codecs {
		if c == nil {
			continue
		}

		m[c.ContentType(nil)] = c
	}

	return m
}
