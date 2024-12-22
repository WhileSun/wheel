package gserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GserverTLS struct {
	CertFile string
	KeyFile  string
}

type GserverConf struct {
	Name string `default:"go-server"`
	Host string `default:"0.0.0.0"`
	Port string `default:"3000"`
	Tls  GserverTLS
}

var (
	httpServer *http.Server
)

func New(gserverConf GserverConf) *GserverConf {
	// 如果没有设置httpServer，创建一个默认的httpServer
	httpServer = &http.Server{
		Addr:         gserverConf.Host + ":" + gserverConf.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &gserverConf
}

func (gserverConf *GserverConf) SetHttpServer(server *http.Server) {
	httpServer = server
}

func (gserverConf *GserverConf) SetHttpHandler(handler http.Handler) {
	httpServer.Handler = handler
}

func (gserverConf *GserverConf) Run() {
	// 创建一个可取消的context.Context用于控制服务器生命周期
	_, cancel := context.WithCancel(context.Background())
	fmt.Printf("[%s] Listening and serving %s \n", gserverConf.Name, httpServer.Addr)
	// 在一个单独的协程中启动HTTP服务器
	go func() {
		if gserverConf.Tls.CertFile != "" && gserverConf.Tls.KeyFile != "" {
			if err := httpServer.ListenAndServeTLS(gserverConf.Tls.CertFile, gserverConf.Tls.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("[%s] HTTPS server start failed: %v", gserverConf.Name, err)
			}
		} else {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("[%s] HTTP server start failed: %v", gserverConf.Name, err)
			}
		}
	}()
	// 创建一个通道用于接收系统信号
	sigChan := make(chan os.Signal, 1)
	// 通知通道接收SIGINT（Ctrl+C）和SIGTERM（终止信号）信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待接收到信号
	<-sigChan
	fmt.Printf("[%s] Received signal, shutting down server... \n", gserverConf.Name)
	// 取消context，触发关闭服务器的逻辑
	cancel()

	// 设置一个超时时间，用于等待服务器优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用Shutdown方法优雅关闭服务器，它会等待正在处理的请求完成后关闭
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("[%s] HTTP server shutdown failed: %v \n", gserverConf.Name, err)
	}
	fmt.Printf("[%s] HTTP server stopped gracefully \n", gserverConf.Name)
}
