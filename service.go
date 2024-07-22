package main

import (
	"context"
	"coursera/hw7_microservice/acl"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// !!!! тут вы пишете код
// !!!! обращаю ваше внимание - в этом задании запрещены глобальные переменные

type (
	MyMicroService struct {
		mu           sync.RWMutex
		ctx          context.Context
		addr         string
		aclData      *acl.ACLS
		server       *grpc.Server
		portListener net.Listener
	}
)

func (s *MyMicroService) stopGRPC() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.server.GracefulStop()
}

func (s *MyMicroService) obtainPort() error {
	var err error
	if s.portListener, err = net.Listen("tcp", s.addr); err != nil {
		fmt.Printf("cant listen port: %s\n", err.Error())
		return err
	}
	return nil
}

func (s *MyMicroService) startGRPC() error {
	var err error

	s.mu.Lock()
	defer s.mu.Unlock()

	if err = s.obtainPort(); err != nil {
		return err
	}

	go func() { // Start grpc service
		s.server = grpc.NewServer()
		// register new GRPC service to bind to port
		RegisterBizServer(s.server, NewBizServer(s.aclData))
		RegisterAdminServer(s.server, NewAdminServer(s.aclData))
		if err := s.server.Serve(s.portListener); err != nil {
			panic(fmt.Errorf("failed to start grpc service : %s", err.Error()))
		}
		fmt.Printf("server started on %s\n", s.addr)
	}()

	go func() { // Wait for any cancelation of grpc service.
		<-s.ctx.Done()
		fmt.Println("grpc service canceled.")
		s.stopGRPC()
	}()
	return nil
}

func StartMyMicroservice(ctx context.Context, addr string, aclData string) error {
	acls, err := acl.BuildACL(aclData)
	if err != nil {
		return err
	}
	serviceInstance := MyMicroService{ctx: ctx, addr: addr, aclData: acls}
	return serviceInstance.startGRPC()
}
