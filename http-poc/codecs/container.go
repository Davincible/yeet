package codecs

import (
	"errors"

	"http-poc/codecs/form"

	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Codec runtime.Marshaler

type Codecs map[string]Codec

type RegCodec struct {
	// One or more content types for which the codec is responsible
	ContentTypes []string
	Codec        Codec
}

var (
	CodecRegistry   = make(Codecs)
	DefaultCodecSet = wire.NewSet(ProvideCodecJSON, ProvideCodecProto, ProvideCodecForm)
	ErrNotProtoType = errors.New("provided interface is not of type proto.Message")
)

func ProvideCodecJSON() RegCodec {
	return RegCodec{
		ContentTypes: []string{
			"application/json",
		},
		Codec: new(runtime.JSONPb),
	}
}

func ProvideCodecProto() RegCodec {
	return RegCodec{
		ContentTypes: []string{
			"application/octet-stream",
			"application/protobuf",
			"application/x-protobuf",
		},
		Codec: new(runtime.ProtoMarshaller),
	}
}

func ProvideCodecForm() RegCodec {
	return RegCodec{
		ContentTypes: []string{
			"x-www-form-urlencoded",
		},
		Codec: form.NewFormCodec(),
	}
}
