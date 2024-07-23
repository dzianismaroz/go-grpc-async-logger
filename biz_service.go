package main

import (
	"context"
)

// BizServer - is just a stub. No actual logic.
func (s *MyMicroService) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *MyMicroService) Add(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil //dummy
}

func (s *MyMicroService) Test(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil // dummy
}
