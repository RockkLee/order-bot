package grpcsync

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Receiver interface {
	ReceiveOrder(ctx context.Context, req SubmitOrderRequest) error
}

type Server struct {
	addr     string
	receiver Receiver
	server   *grpc.Server
}

func NewServer(addr string, receiver Receiver) *Server {
	codec := JSONCodec{}
	gs := grpc.NewServer(grpc.ForceServerCodec(codec))
	s := &Server{addr: addr, receiver: receiver, server: gs}
	gs.RegisterService(&grpc.ServiceDesc{
		ServiceName: "ordersync.OrderSyncService",
		HandlerType: (*Receiver)(nil),
		Methods: []grpc.MethodDesc{{
			MethodName: "SubmitOrder",
			Handler:    s.submitOrder,
		}},
	}, receiver)
	return s
}

func (s *Server) submitOrder(_ any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req SubmitOrderRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	if err := s.receiver.ReceiveOrder(ctx, req); err != nil {
		return nil, err
	}
	return &SubmitOrderResponse{Accepted: true}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpcsync.Server.Start: %w", err)
	}
	go func() { _ = s.server.Serve(lis) }()
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
