package go_web

import (
	"net/http"
	"time"

	"github.com/cyuzhe1994-commits/go-web/middleware"
)

// 框架核心结构
type Engine struct {
	logger IFrameWorkLog
	Router *Router
}

func NewEngine(logger IFrameWorkLog) *Engine {
	if logger == nil {
		logger = NewDefaultFrameWorkLog(DefaultFrameWorkLogLevelInfo, time.UTC)
	}
	router := NewRouter()
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)
	return &Engine{logger: logger, Router: router}
}

func (e *Engine) Run(addr string) (err error) {
	e.logger.Info("Starting server at %s", addr)
	// 这里可以添加启动服务器的逻辑，例如监听端口等
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.logger.Info("收到来自 %s 的请求: %s", r.RemoteAddr, r.URL.Path)
	// 这里可以添加处理 HTTP 请求的逻辑，例如路由分发等
	e.Router.Handle(w, r)
}
