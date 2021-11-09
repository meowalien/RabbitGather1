package main

import (
	_ "core/src/init"

	"context"
	"core/src/module/db/mariadb"
	"core/src/module/db/redisdb"
	"core/src/module/log"
	"core/src/server"
	"core/src/service/member"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


func finalize() {
	e := redisdb.Conn.Close()
	if e != nil {
		fmt.Println("error when cosing redisdb.Conn: ", e.Error())
		e = nil
	}
	e = mariadb.Conn.Close()
	if e != nil {
		fmt.Println("error when cosing mariadb.GORMConn(mysql): ", e.Error())
		e = nil
	}}


func main() {
	defer finalize()
	fmt.Println("main ...")
	log.Logger.Debug("Logger.Debug test ")
	log.Logger.Info("Logger.Info test")
	log.Logger.Warn("Logger.Warn test")
	log.Logger.Error("Logger.Error test")
	ctx := context.Background()


	memberService := member.Member{}
	err := memberService.Initialize(ctx)
	if err != nil {
		panic(err.Error())
	}




	httpServer := server.InitGinHTTPServer(ctx,memberService)
	httpServer.Start(ctx)

	rpcServer := server.InitGRPCServer(ctx)
	rpcServer.Start(ctx)

	/* 等待結束命令 */

	waitForShutdown(ctx, func() error {
		e := rpcServer.Stop(ctx)
		if e != nil {
			fmt.Println("error when close GRPCServer")
		}

		e1 := httpServer.Stop(ctx)
		if e1 != nil {
			fmt.Println("error when close HTTPServer")
		}
		return nil
	})

}

func waitForShutdown(ctx context.Context, callbackFunc func() error) {
	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	// 所有服務都停止了
	case <-server.NoServerUp:
		fmt.Println("No active server running, ending main function")
	case <-ctx.Done():
		fmt.Println("Shutdown with Context done")
	case <-quitSignal:
		fmt.Println("Shutdown with OS QuitSignal")
	}

	err := callbackFunc()
	if err != nil {
		fmt.Println("Shutdown Error: ", err.Error())
	}
}
