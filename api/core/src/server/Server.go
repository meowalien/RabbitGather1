package server

import (
	"context"
	"sync/atomic"
)

var _TotalUpHTTPServer int32 = 0

// NoServerUp will be close when no server alive
var NoServerUp = make(chan interface{}, 1)

func AddTotalUpServer(i int32) int32 {
	total := atomic.AddInt32(&_TotalUpHTTPServer, i)
	if total <= 0 {
		close(NoServerUp)
		return 0
	}
	return total
}

// 有幾個伺服器是上線的
func GetTotalUpServer() int32 {
	return _TotalUpHTTPServer
}

type Server interface {
	Startable
	Stopable
}

type Startable interface {
	Start(ctx context.Context)
}
type Stopable interface {
	Stop(ctx context.Context) error
}
