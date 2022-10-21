package handler

import (
	"context"
	"fmt"

	"http-poc/proto"
)

var _ proto.StreamsServer = (*EchoHandler)(nil)

type EchoHandler struct {
	proto.UnimplementedStreamsServer
}

func (c *EchoHandler) Call(ctx context.Context, in *proto.CallRequest) (*proto.CallResponse, error) {
	fmt.Println("Received message from client:", in.Name)
	return &proto.CallResponse{Msg: "Message received!"}, nil
}
