package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Interceptor for streamed endpoints of Admins servrice : logs and statistics.
// Provides authentication and broadcasting statistics and logging.
func (s *MyMicroService) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ssCtx := ss.Context()
	md, _ := metadata.FromIncomingContext(ssCtx)
	acls := s.acl
	principal := md.Get(consumerKey)
	// ----- Authentication.
	if authErr := acls.Authenticate(principal, info.FullMethod); authErr != nil {
		return authErr
	}
	// ---- broadcasting statistics and logging to subscribers.
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

// Unary interceptor for Biz service endpoints.
// Provides authentication and broadcasting statistics and logging.
func (s *MyMicroService) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	reply, err := handler(ctx, req)
	md, _ := metadata.FromIncomingContext(ctx)
	principal := md.Get(consumerKey)
	acls := s.acl
	// ----- Authentication.
	if authErr := acls.Authenticate(principal, info.FullMethod); authErr != nil {
		return nil, authErr
	}
	// ---- broadcasting statistics and logging to subscribers.
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
