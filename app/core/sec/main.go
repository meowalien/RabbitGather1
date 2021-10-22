package main

import (
	"context"
	"core/sec/conf"
	"core/sec/db/mariadb"
	"core/sec/db/redisdb"
	"core/sec/lib/config"
	"core/sec/log"
	"core/sec/module/permission"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	sys_log "log"
	"os"
	"os/signal"
	"syscall"
)

func InitFlags() {
	fmt.Println("InitFlags ...")
	flag.BoolVar(&conf.DEBUG_MOD, "debug", false, "\"true\" to open debug mode")

	flag.Parse()

	fmt.Println("debug mod: ", conf.DEBUG_MOD)
}

const config_file = "config/config.json"

func InitConfig() {
	fmt.Println("InitConfig ...")
	err := config.JsonConfigModleMapping(&conf.GlobalConfig, config_file)
	if err != nil {
		panic(err.Error())
	}
	_, e := pretty.Println(conf.GlobalConfig)
	if e != nil {
		fmt.Println("error when printing GlobalConfig by pretty:", err.Error())
		fmt.Println(conf.GlobalConfig)
	}
}

func init() {
	InitFlags()
	InitConfig()

	redisdb.InitRedis()
	addToFinalize(func() {
		e := redisdb.Conn.Close()
		if e != nil {
			fmt.Println("error when cosing redisdb.Conn: ", e.Error())
		}
	})

	mariadb.InitMariadbDBConnection()

	// 映攝所需資料結構
	//module.MappingModules()


	addToFinalize(func() {
		e := mariadb.Conn.Close()
		if e != nil {
			fmt.Println("error when cosing mariadb.GORMConn(mysql): ", e.Error())
		}
	})

	permission.InitRBAC(mariadb.Conn)

	log.InitLogger()
}

func main() {
	defer finalize()

	/* 啟動服務單元 */

	memberService := member.MemberService{}
	err := memberService.Initialize()
	if err != nil {
		panic(err.Error())
	}

	// 大廳服務
	lobbyService := lobby.LobbyService{}
	err = lobbyService.Initialize()
	if err != nil {
		panic(err.Error())
	}

	// 玩家錢包操作
	walletService := wallet.WalletService{}
	err = walletService.Initialize()
	if err != nil {
		panic(err.Error())
	}

	// 第三方金流
	thirdPartyCashFlowService := third_party_cash_flow.ThirdPartyCashFlowService{}
	err = thirdPartyCashFlowService.Initialize()
	if err != nil {
		panic(err.Error())
	}
	/* 啟動伺服器，並掛載服務*/

	ctx := context.Background()

	httpServer := server.InitGinHTTPServer(ctx, memberService, lobbyService, thirdPartyCashFlowService)
	httpServer.Start(ctx)

	rpcServer := server.InitGRPCServer(ctx, memberService, walletService)
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

func initForDebug() {
	// 創建測試用帳號
	testUse := user.CreateNewUserRequest{Id: 11652, Name: "U571633915", Role: []string{"login"}, Password: "U571633915", Introducer: "joy-games", InitBalance: 9999999999999}
	fmt.Println("create a testing user : ", testUse)
	_, err := user.CreateNewUser(testUse)
	if err != nil {
		fmt.Println("error when create test user: ", err.Error())
		return
	}
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
		sys_log.Println("Shutdown Error: ", err.Error())
	}
}

var finalizeFuncList = []func(){}

func addToFinalize(f func()) {
	finalizeFuncList = append(finalizeFuncList, f)
}
func finalize() {
	for _, f := range finalizeFuncList {
		f()
	}

	wallet_module.EventBroker.Stop()
}
