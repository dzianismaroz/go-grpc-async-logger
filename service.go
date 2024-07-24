package main

import (
	"context"
	"coursera/hw7_microservice/acl"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

// !!!! this is where you write the code
// !!!! Please note that global variables are not allowed in this task.

const consumerKey = "consumer" // Extract consumer as SID from GRPC context on method authentication.

type (
	AclHolder interface {
		GetACLS() *acl.ACLS
	}
	MyMicroService struct {
		UnimplementedAdminServer                          // Satisfy interface of Admin Server.
		UnimplementedBizServer                            // Satisfy interface of Biz Server.
		ctx                      context.Context          // General Context.
		addr                     string                   // GRPC server address to bind.
		acl                      *acl.ACLS                // ACL resolver. It will authenticate GRPC requests.
		server                   *grpc.Server             // GRPC server instance.
		portListener             net.Listener             // Handy port resolving.
		logListeners             map[chan *Event]struct{} // Broadcast streamed logs subscribers.
		statListeners            map[chan *Stat]struct{}  // Broadcast streamed statistic subscribers.
		broadcastStatCh          chan *Stat               // Broadcast notification of all stat subscribers.
		addStatListenerCh        chan chan *Stat          // Add new statistic suscriber.
		removeStatListenerCh     chan chan *Stat          // Unsusrcribe statistic listener.
		broadcastLogCh           chan *Event              // Broadcast notification of all log listeners.
		addLogListenerCh         chan chan *Event         // Add new log suscriber.
		removeLogListenerCh      chan chan *Event         // Remove log subscriber.
	}
)

func newMicroservice(ctx context.Context, addr string, acls *acl.ACLS) *MyMicroService {
	return &MyMicroService{
		ctx:                  ctx,                        // General context.
		addr:                 addr,                       // GRPC server address.
		acl:                  acls,                       // ACL to resolve permissions by its own.
		logListeners:         map[chan *Event]struct{}{}, // Broadcast streamed log listeners.
		statListeners:        map[chan *Stat]struct{}{},  // Broadcast stream statistics listeners.
		broadcastLogCh:       make(chan *Event),
		addLogListenerCh:     make(chan chan *Event),
		broadcastStatCh:      make(chan *Stat),
		addStatListenerCh:    make(chan chan *Stat),
		removeStatListenerCh: make(chan chan *Stat),
		removeLogListenerCh:  make(chan chan *Event),
	}
}

// Authenticate with internal ACL instance.
func (s *MyMicroService) Authenticate(sid []string, path string) error {
	return s.acl.Authenticate(sid, path)
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

	// Add endpoints interceptors for authentication and broadcasting.
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
	}()
	time.Sleep(20 * time.Millisecond)
	go func() { // Wait for any cancelation of grpc service.
		<-s.ctx.Done()
		s.server.GracefulStop()
	}()
	time.Sleep(20 * time.Millisecond)

	// Handle logs and statistics subscription and broadcasting.
	// Must be fired after general server start.
	s.startBroadcasting()
	return nil
}

func StartMyMicroservice(ctx context.Context, addr string, aclData string) error {
	acls, err := acl.BuildACL(aclData)
	if err != nil { // Fail on any problem parsing ACL data provided.
		return err
	}
	return newMicroservice(ctx, addr, acls).startGRPC()
}
