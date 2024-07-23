package main

import (
	"context"
	"coursera/hw7_microservice/acl"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// !!!! this is where you write the code
// !!!! Please note that global variables are not allowed in this task.

const (
	consumer          = "consumer"
	adminMethodPrefix = "/main.Admin/"
	bizMethodPrefix   = "/main.Biz/"
)

type (
	AclHolder interface {
		GetACLS() *acl.ACLS
	}
	MyMicroService struct {
		UnimplementedAdminServer
		UnimplementedBizServer
		ctx          context.Context
		addr         string
		acl          *acl.ACLS
		server       *grpc.Server
		portListener net.Listener

		logListeners  []chan *Event
		statListeners []chan *Stat

		broadcastStatCh   chan *Stat
		addStatListenerCh chan chan *Stat

		broadcastLogCh   chan *Event
		addLogListenerCh chan chan *Event
	}
)

func newMicroservice(ctx context.Context, addr string, acls *acl.ACLS) *MyMicroService {
	return &MyMicroService{
		ctx:               ctx,
		addr:              addr,
		acl:               acls,
		logListeners:      []chan *Event{},
		statListeners:     []chan *Stat{},
		broadcastLogCh:    make(chan *Event),
		addLogListenerCh:  make(chan chan *Event),
		broadcastStatCh:   make(chan *Stat),
		addStatListenerCh: make(chan chan *Stat),
	}
}

func (s *MyMicroService) Authenticate(sid []string, path string) error {
	return s.acl.Authenticate(sid, path)
}

func (s *MyMicroService) stopGRPC() {
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

	if err = s.obtainPort(); err != nil {
		return err
	}
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryInterceptor),
		grpc.StreamInterceptor(s.streamInterceptor),
	)
	// Register services.
	RegisterBizServer(s.server, s)
	RegisterAdminServer(s.server, s)

	go func() { // Start grpc service

		if err := s.server.Serve(s.portListener); err != nil {
			panic(fmt.Errorf("failed to start grpc service : %s", err.Error()))
		}
		fmt.Printf("server started on %s\n", s.addr)
	}()

	go func() {
		for {
			select {
			case ch := <-s.addLogListenerCh:
				s.logListeners = append(s.logListeners, ch)
			case event := <-s.broadcastLogCh:
				for _, ch := range s.logListeners {
					ch <- event
				}
			case <-s.ctx.Done():
				s.stopGRPC()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case ch := <-s.addStatListenerCh:
				s.statListeners = append(s.statListeners, ch)
			case stat := <-s.broadcastStatCh:
				for _, ch := range s.statListeners {
					ch <- stat
				}
			case <-s.ctx.Done():
				return
			}
		}
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
	if err != nil { // Fail on any problem parsing ACL data provided.
		return err
	}
	return newMicroservice(ctx, addr, acls).startGRPC()
}
