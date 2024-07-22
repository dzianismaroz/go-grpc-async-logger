package main

import "coursera/hw7_microservice/acl"

type MyAdminServer struct {
	UnimplementedAdminServer
	acls *acl.ACLS
}

func NewAdminServer(acls *acl.ACLS) AdminServer {
	return &MyAdminServer{acls: acls}
}

func (s *MyMicroService) Logging(*Nothing, Admin_LoggingServer) error {

	return nil
}

func (s *MyMicroService) Statistics(*StatInterval, Admin_StatisticsServer) error {

	return nil
}
