package server

import "context"

type HelloService struct{}

func (h *HelloService) SayHelloWorld(ctx context.Context,r )
func NewHelloService() *HelloService {
	return &HelloService{}
}
