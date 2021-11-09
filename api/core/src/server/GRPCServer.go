package server

import (
	"context"
	"core/src/conf"
	"core/src/module/log"
	"core/src/service"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func InitGRPCServer(ctx context.Context, registers ...service.GRPCServerRegister) *GRPCServer {

	firstServer := GRPCServer{
		RPCServerConfig: conf.GlobalConfig.Server.RPCServer,
		MountService: func(ctx context.Context, rpcServer *grpc.Server) error {
			for _, register := range registers {
				err := register.RegisterRPCServer(rpcServer)
				if err != nil {
					return fmt.Errorf("error when start memberService-RPC: %w", err)
				}
			}
			return nil
		},
	}
	return &firstServer
}

type GRPCServer struct {
	ID int32
	conf.RPCServerConfig
	BeforeStop   func(ctx context.Context)
	MountService func(ctx context.Context, rpcServer *grpc.Server) error
	rpcServer    *grpc.Server
}

func (r *GRPCServer) Start(ctx context.Context) {
	if r.MountService == nil {
		panic("MountService is required")
	}
	if r.Port == "" {
		panic("Port is required")
	}
	if r.Host == "" {
		panic("Port is required")
	}

	lis, err := net.Listen("tcp", r.Host+":"+r.Port)
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	s := grpc.NewServer()

	err = r.MountService(ctx, s)

	if err != nil {
		log.Logger.Error("fail to start GRPCServer, MountService error: ", err.Error())
		return
	}
	r.rpcServer = s

	r.ID = AddTotalUpServer(1)

	if r.Name == "" {
		r.Name = fmt.Sprintf("GRPCServer ID:%d", r.ID)
	} else {
		r.Name = fmt.Sprintf("GRPCServer ID:%d %s", r.ID, r.Name)
	}

	var hasErr = false
	go func() {
		if err := s.Serve(lis); err != nil {
			hasErr = true
			log.Logger.Errorf("error when start %s err: %s", r.Name, err.Error())
		}
		AddTotalUpServer(-1)
	}()
	if hasErr {
		return
	}

	log.Logger.Infof("RPC Server: %s listening at: %v", r.Name,r.Host+":"+r.Port)
}

func (r *GRPCServer) Stop(ctx context.Context) error {
	if r.BeforeStop != nil {
		r.BeforeStop(ctx)
	}

	r.rpcServer.GracefulStop()
	log.Logger.Infof("%s closed.", r.Name)

	return nil
}
