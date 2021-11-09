package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// 向Gin Router伺服器掛載服務
type Service interface {
	Initialize(ctx context.Context) error
}

// 向Gin Router伺服器掛載服務
type GinServerRegister interface {
	//HTTPServerRegister
	// 掛載服務提供的處理器
	Mount(ctx context.Context, router *gin.Engine) error
}

// 向GRPC伺服器掛載服務
type GRPCServerRegister interface {
	//RPCServerRegister
	RegisterRPCServer(rpcServer grpc.ServiceRegistrar) error
}
