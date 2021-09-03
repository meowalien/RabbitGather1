package main

import (
	"context"
	"fmt"
	"github.com/meowalien/go-util"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)
func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "../config/web_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	//log.DEBUG.Println("WebServer - ServePath : ", ServePath)
	if err != nil {
		panic(err.Error())
	}
}
type WebServer struct {

}

func (s WebServer) Startup(ctx context.Context) error {

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	webserver := WebServer{}
	err := webserver.Startup(ctx)
	if err != nil {
		cancel()
		panic(err.Error())
	}

	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		fmt.Println("Shutdown with Context done")
	case <-quitSignal:
		cancel()
		fmt.Println("Shutdown with OS QuitSignal")
	}
}
