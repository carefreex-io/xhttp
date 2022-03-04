package xhttp

import (
	"context"
	"github.com/carefreex-io/config"
	"github.com/carefreex-io/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type XHttp struct {
	Engine      *gin.Engine
	server      *http.Server
	onStartFunc []func()
	onStopFunc  []func()
}

func NewXHttp() *XHttp {
	return &XHttp{
		Engine: gin.New(),
	}
}

func (h *XHttp) RegisterGlobalMiddleware(middles []func(c *gin.Context)) {
	for _, middle := range middles {
		h.Engine.Use(middle)
	}
}

func (h *XHttp) AddOnStartFunc(fn []func()) {
	h.onStartFunc = append(h.onStartFunc, fn...)
}

func (h *XHttp) AddOnStopFunc(fn []func()) {
	h.onStopFunc = append(h.onStopFunc, fn...)
}

func (h *XHttp) Start() {
	h.onStart()

	h.server = &http.Server{
		Addr:           ":" + config.GetString("Service.Port"),
		Handler:        h.Engine,
		ReadTimeout:    config.GetDuration("HttpServer.ReadTimeout") * time.Second,
		WriteTimeout:   config.GetDuration("HttpServer.WriteTimeout") * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panicf("server start failed: err=%v", err)
			return
		}
	}()

	h.shutdown()
}

func (h *XHttp) onStart() {
	for _, fn := range h.onStartFunc {
		fn()
	}
}

func (h *XHttp) onStop() {
	for _, fn := range h.onStopFunc {
		fn()
	}
}

func (h *XHttp) shutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown Server ...")

	h.onStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		logger.Println("timeout of 5 seconds.")
	}
	logger.Println("Server exiting")
}
