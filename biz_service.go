package main

import (
	"context"
	"coursera/hw7_microservice/acl"
	"fmt"
)

type MyBizServer struct {
	UnimplementedBizServer
	acls *acl.ACLS
}

func NewBizServer(acls *acl.ACLS) BizServer {
	return &MyBizServer{acls: acls}
}

// BizServer
func (s *MyBizServer) Check(ctx context.Context, in *Nothing) (*Nothing, error) {

	return nil, nil
}

func (s *MyBizServer) Add(ctx context.Context, in *Nothing) (*Nothing, error) {

	return nil, nil
}

func (s *MyBizServer) Test(ctx context.Context, in *Nothing) (*Nothing, error) {
	fmt.Println("WE ENTERED !!!! #####")
	return nil, nil
}
