package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *MyMicroService) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ssCtx := ss.Context()
	md, _ := metadata.FromIncomingContext(ssCtx)
	acls := s.acl
	principal := md.Get(consumer)

	if authErr := acls.Authenticate(principal, info.FullMethod); authErr != nil {
		return authErr
	}

	s.broadcastLogCh <- &Event{
		Consumer: principal[0],
		Method:   info.FullMethod,
		Host:     "127.0.0.1:8083",
	}
	s.broadcastStatCh <- &Stat{
		ByConsumer: map[string]uint64{principal[0]: 1},
		ByMethod:   map[string]uint64{info.FullMethod: 1},
	}
	return handler(srv, ss)
}

func (s *MyMicroService) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	reply, err := handler(ctx, req)
	md, _ := metadata.FromIncomingContext(ctx)
	principal := md.Get(consumer)
	acls := s.acl
	fmt.Printf("###MD: %#v\n", md)
	if authErr := acls.Authenticate(principal, info.FullMethod); authErr != nil {
		return nil, authErr
	}
	s.broadcastLogCh <- &Event{
		Consumer: principal[0],
		Method:   info.FullMethod,
		Host:     "127.0.0.1:8083",
	}
	s.broadcastStatCh <- &Stat{
		ByConsumer: map[string]uint64{principal[0]: 1},
		ByMethod:   map[string]uint64{info.FullMethod: 1},
	}

	return reply, err
}
