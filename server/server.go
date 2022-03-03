package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"fire-scaffold/conf"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Server *http.Server
}

const (
	defaultTimeout       = 30 * time.Second
	defaulMaxHeaderBytes = 1 << 20
)

func New(mode, addr string, opts ...Option) *HttpServer {
	gin.SetMode(mode)

	r := InitRouter(TimeoutMiddleware)

	new := &HttpServer{Server: &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    defaultTimeout,
		WriteTimeout:   defaultTimeout,
		MaxHeaderBytes: defaulMaxHeaderBytes,
	}}

	for _, opt := range opts {
		if opt == nil {
			panic(" [PANIC] option not to be nil!")
		}

		opt(new)
	}

	return new
}

func (h *HttpServer) Run() {
	go func() {
		log.Printf(" [INFO] HttpServerRun:%s\n", conf.GlobalConfig.HTTP.Addr)

		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", conf.GlobalConfig.HTTP.Addr, err)
		}
	}()
}

func (h *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}

	log.Printf(" [INFO] HttpServerStop stopped\n")
}
