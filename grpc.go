// Package grpc provides a grpc transport
package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	pb "github.com/unistack-org/micro-network-transport-grpc/proto"
	"github.com/unistack-org/micro/v3/network/transport"
	mnet "github.com/unistack-org/micro/v3/util/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type grpcTransport struct {
	opts transport.Options
}

type grpcTransportListener struct {
	listener net.Listener
	secure   bool
	tls      *tls.Config
	opts     transport.ListenOptions
}

func (t *grpcTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *grpcTransportListener) Close() error {
	return t.listener.Close()
}

func (t *grpcTransportListener) Accept(fn func(transport.Socket)) error {
	var opts []grpc.ServerOption

	// setup tls if specified
	if t.secure && t.tls == nil {
		return fmt.Errorf("request secure communication but *tls.Config is nil")
	} else if t.secure && t.tls != nil {
		creds := credentials.NewTLS(t.tls)
		opts = append(opts, grpc.Creds(creds))
	}

	// new service
	srv := grpc.NewServer(opts...)

	// register service
	pb.RegisterTransportServer(srv, &microTransport{addr: t.listener.Addr().String(), fn: fn})

	// start serving
	return srv.Serve(t.listener)
}

func (t *grpcTransport) Dial(ctx context.Context, addr string, opts ...transport.DialOption) (transport.Client, error) {
	dopts := transport.NewDialOptions(opts...)

	options := []grpc.DialOption{}

	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		creds := credentials.NewTLS(config)
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	// dial the server
	ctx, cancel := context.WithTimeout(ctx, dopts.Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, options...)
	if err != nil {
		return nil, err
	}

	// create stream
	stream, err := pb.NewTransportClient(conn).Stream(context.Background())
	if err != nil {
		return nil, err
	}

	// return a client
	return &grpcTransportClient{
		conn:   conn,
		stream: stream,
		local:  "localhost",
		remote: addr,
	}, nil
}

func (t *grpcTransport) Listen(ctx context.Context, addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	options := transport.NewListenOptions(opts...)

	ln, err := mnet.Listen(addr, func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	})
	if err != nil {
		return nil, err
	}

	return &grpcTransportListener{
		listener: ln,
		tls:      t.opts.TLSConfig,
		secure:   t.opts.Secure,
		opts:     options,
	}, nil
}

func (t *grpcTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *grpcTransport) Options() transport.Options {
	return t.opts
}

func (t *grpcTransport) String() string {
	return "grpc"
}

func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &grpcTransport{opts: options}
}
