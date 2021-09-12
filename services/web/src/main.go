package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meowalien/go-util/json"
	"github.com/meowalien/go-util/logger"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)


var ServePath *url.URL
func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := json.ParseFileJsonConfig(&config, "config/web_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	if err != nil {
		panic(err.Error())
	}
}
type WebServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var log *logger.LoggerWrapper

func init() {
	//logger.LogLevelMask  = logger.ALL
	log  = logger.NewLoggerWrapper("WebServer","log/")
}

func (s *WebServer) Startup(ctx context.Context) error {
	addr :=   ":" + ServePath.Port()
	log.DEBUG.Printf("WebServer listen on : \"%s\"\n",addr)

	s.ginEngine = gin.Default()
	s.serverInst = &http.Server{
		Addr:   addr,
		Handler: s.ginEngine,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}

	s.MountService(ctx)
	go func() {
		if err := s.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ERROR.Println(err.Error())
		}
	}()
	log.DEBUG.Println("WebServer Started .")
	return nil
}

func (s *WebServer) MountService(ctx context.Context) {
	s.ginEngine.NoRoute(func(c *gin.Context) {
		c.File("assets/index.html")
	})
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
