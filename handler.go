package grpc

import (
	"runtime/debug"

	pb "go.unistack.org/micro-network-transport-grpc/v3/proto"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/network/transport"
	"google.golang.org/grpc/peer"
)

// microTransport satisfies the pb.TransportServer inteface
type microTransport struct {
	pb.UnimplementedTransportServer
	addr string
	fn   func(transport.Socket)
}

func (m *microTransport) Stream(ts pb.Transport_StreamServer) (err error) {
	sock := &grpcTransportSocket{
		stream: ts,
		local:  m.addr,
	}

	p, ok := peer.FromContext(ts.Context())
	if ok {
		sock.remote = p.Addr.String()
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error(ts.Context(), "panic recovered: ", r)
			logger.Error(ts.Context(), string(debug.Stack()))
			sock.Close()
			err = errors.InternalServerError("go.micro.transport", "panic recovered: %v", r)
		}
	}()

	// execute socket func
	m.fn(sock)

	return err
}
