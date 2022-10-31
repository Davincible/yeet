package handler

import (
	"context"

	"http-poc/proto"
)

var _ proto.StreamsServer = (*EchoHandler)(nil)

type EchoHandler struct {
	proto.UnimplementedStreamsServer
}

func (c *EchoHandler) Call(ctx context.Context, in *proto.CallRequest) (*proto.CallResponse, error) {
	return &proto.CallResponse{Msg: "Hello " + in.Name}, nil
}
