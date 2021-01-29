// Code generated by protoc-gen-micro
// source: transport.proto
package transport

import (
	"context"

	micro_client "github.com/unistack-org/micro/v3/client"
	micro_server "github.com/unistack-org/micro/v3/server"
)

var (
	_ micro_server.Option
	_ micro_client.Option
)

type transportService struct {
	c    micro_client.Client
	name string
}

// Micro client stuff

// NewTransportService create new service client
func NewTransportService(name string, c micro_client.Client) TransportService {
	return &transportService{c: c, name: name}
}

func (c *transportService) Stream(ctx context.Context, opts ...micro_client.CallOption) (Transport_StreamService, error) {
	stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "Transport.Stream", &Message{}), opts...)
	if err != nil {
		return nil, err
	}
	return &transportServiceStream{stream}, nil
}

type transportServiceStream struct {
	stream micro_client.Stream
}

func (x *transportServiceStream) Close() error {
	return x.stream.Close()
}

func (x *transportServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *transportServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *transportServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *transportServiceStream) Send(m *Message) error {
	return x.stream.Send(m)
}

func (x *transportServiceStream) Recv() (*Message, error) {
	m := &Message{}
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
} // Micro server stuff

type transportHandler struct {
	TransportHandler
}

func (h *transportHandler) Stream(ctx context.Context, stream micro_server.Stream) error {
	return h.TransportHandler.Stream(ctx, &transportStreamStream{stream})
}

type transportStreamStream struct {
	stream micro_server.Stream
}

func (x *transportStreamStream) Close() error {
	return x.stream.Close()
}

func (x *transportStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *transportStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *transportStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}
func (x *transportStreamStream) Send(m *Message) error {
	return x.stream.Send(m)
}

func (x *transportStreamStream) Recv() (*Message, error) {
	m := &Message{}
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
