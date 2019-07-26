package dispatcher

import (
	"context"
	"google.golang.org/grpc"
)

type (
	NewClientFunc func(*grpc.ClientConn) interface{}
	WorkerFunc    func(interface{}) error
)

type Service struct {
	Ctx  context.Context
	Call func(*grpc.ClientConn) interface{}
}
