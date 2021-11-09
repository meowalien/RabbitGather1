package server

import (
	"context"
	"core/src/conf"
	"core/src/module/log"
	"core/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func CrossHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//origin := c.Request.Header.Get("Origin")
		//if origin == "" {
		//	c.AbortWithStatus(http.StatusForbidden)
		//}
		//ok := CheckOrigin(origin)
		//if !ok {
		//	c.AbortWithStatus(http.StatusForbidden)
		//	return
		//}

		//接收客戶端傳送的origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		//伺服器支援的所有跨域請求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		//允許跨域設定可以返回其他子段，可以自定義欄位
		//c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, "+
		//	"X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, "+
		//	"X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, "+
		//	"Content-Type, Pragma, token, openid, opentoken, Authentication-Token")
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		//允許瀏覽器(客戶端)可以解析的頭部 (重要)
		c.Header("Access-Control-Expose-Headers", exposeHeaders)
		//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
		//	"Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, "+
		//	"Expires, Last-Modified, Pragma, FooBar")
		//設定快取時間
		c.Header("Access-Control-Max-Age", "86400")
		//允許客戶端傳遞校驗資訊比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "false")

		//允許型別校驗
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			//c.JSON(http.StatusOK, "ok!")
		}

		c.Next()
	}
}

var allowOrigin string

func init() {
	allowOrigin = strings.Join(conf.GlobalConfig.Server.AllowOrigin, ",")
}
var allowHeaders string


var exposeHeaders string

func init() {
	exposeHeaders = strings.Join(conf.GlobalConfig.Server.ExposeHeaders, ",")
}


func InitGinHTTPServer(ctx context.Context, registers ...service.GinServerRegister) *GinHTTPServer {
	server := GinHTTPServer{
		HTTPServerConfig: conf.GlobalConfig.Server.HTTPServer,
		MountService: func(ctx context.Context, router *gin.Engine) error {
			router.Use(gin.Recovery())

			router.Use(CrossHandler())

			if registers == nil {
				return nil
			}
			for _, register := range registers {
				err := register.Mount(ctx, router)
				if err != nil {
					return fmt.Errorf("error when start memberService-RPC: %w", err)
				}
			}
			return nil
		},
	}
	return &server
}

type GinHTTPServer struct {
	ID           int32
	Debug        bool
	MountService func(ctx context.Context, router *gin.Engine) error
	BeforeStop   func(ctx context.Context)
	//ShutdownWaitTime time.Duration
	conf.HTTPServerConfig
	serverInst              *http.Server
	ginEngine               *gin.Engine
	shutdownCallbackMethods []func() error
}

func (g *GinHTTPServer) Start(ctx context.Context) {
	if g.MountService == nil {
		panic("MountService is required")
	}
	if g.Port == "" {
		panic("Port is required")
	}
	if g.Host == "" {
		panic("Port is required")
	}

	server := gin.Default()

	err := g.MountService(ctx, server)
	if err != nil {
		log.Logger.Error("fail to start GinHTTPServer, MountService error: ", err.Error())
		return
	}

	g.serverInst = &http.Server{
		Addr:    g.Host + ":" + g.Port,
		Handler: server,
	}

	g.ID = AddTotalUpServer(1)

	if g.Name == "" {
		g.Name = fmt.Sprintf("HTTPServer UUID:%d", g.ID)
	} else {
		g.Name = fmt.Sprintf("HTTPServer UUID:%d %s", g.ID, g.Name)
	}

	var hasErr = false
	go func() {
		if err := g.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hasErr = true
			log.Logger.Errorf("error when start %s err: %s", g.Name, err.Error())
		}
		AddTotalUpServer(-1)
	}()

	if hasErr {
		return
	}
	log.Logger.Infof("HTTP Server: %s listening at: %s", g.Name, g.Host+":"+g.Port)
}

var DefaultShutdownWaitTime = time.Minute * 3

func (g *GinHTTPServer) Stop(ctx context.Context) error {
	if g.ShutdownWaitTime == 0 {
		g.ShutdownWaitTime = DefaultShutdownWaitTime
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.ShutdownWaitTime)
	defer cancel()

	if g.BeforeStop != nil {
		g.BeforeStop(ctx)
	}
	if err := g.serverInst.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown %s, err: %w", g.Name, err)
	} else {
		log.Logger.Infof("%s closed", g.Name)
		return nil
	}
}
