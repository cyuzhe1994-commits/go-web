package go_web

import (
	"net/http"
	"time"

	"github.com/cyuzhe1994-commits/go-web/middleware"
	"github.com/cyuzhe1994-commits/go-web/public"
)

// 框架核心结构
type Engine struct {
	logger      IFrameWorkLog
	Router      *Router
	middlewares []public.Middleware
}

func NewEngine(logger IFrameWorkLog) *Engine {
	if logger == nil {
		logger = NewDefaultFrameWorkLog(DefaultFrameWorkLogLevelInfo, time.UTC)
	}
	router := NewRouter()
	middlewares := make([]public.Middleware, 0)
	middlewares = append(middlewares, middleware.Recovery)
	middlewares = append(middlewares, middleware.Logger)
	return &Engine{logger: logger, Router: router}
}

func (e *Engine) Use(middleware public.Middleware) {
	e.middlewares = append(e.middlewares, middleware)
}

func (e *Engine) Run(addr string) (err error) {
	e.logger.Info("Starting server at %s", addr)
	// 这里可以添加启动服务器的逻辑，例如监听端口等
	return http.ListenAndServe(addr, e)
}

func (e *Engine) combineMiddlewares(handler public.HandlerFunc) public.HandlerFunc {
	// 从后往前执行，这样最先 Use 的中间件会在最外层执行
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		handler = e.middlewares[i](handler)
	}
	return handler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.logger.Info("收到来自 %s 的请求: %s", r.RemoteAddr, r.URL.Path)
	// 这里可以添加处理 HTTP 请求的逻辑，例如路由分发等
	handler, ctx := e.Router.Handle(w, r)
	if handler == nil {
		handler = func(ctx *public.Context) {
			ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "404 not found",
			})
		}
	}
	finalHandler := e.combineMiddlewares(handler)
	if ctx == nil {
		ctx = &public.Context{
			Request: r,
			Writer:  w,
		}
	}
	finalHandler(ctx)
}
