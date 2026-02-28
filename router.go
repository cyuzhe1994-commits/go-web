package go_web

import (
	"net/http"
	"strings"

	"github.com/cyuzhe1994-commits/go-web/route"
)

type Router struct {
	trees      map[string]*route.Tree      // 在新建 group 后，树结构是共享的
	handlers   map[*route.Node]HandlerFunc // 在新建 group 后，handlers 是共享的，保证每个 group 可以独立添加路由
	middleware []Middleware                // 在新建 group 后，中间件是复制的，保证每个 group 可以独立添加中间件
	prefix     string
}

func NewRouter() *Router {
	return &Router{
		trees:      make(map[string]*route.Tree),
		middleware: make([]Middleware, 0),
		prefix:     "",
		handlers:   make(map[*route.Node]HandlerFunc),
	}
}

func (r *Router) Use(middleware Middleware) {
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) Group(prefix string) *Router {

	newMiddleware := make([]Middleware, len(r.middleware))
	copy(newMiddleware, r.middleware)

	return &Router{
		trees:      r.trees,
		middleware: newMiddleware,
		handlers:   make(map[*route.Node]HandlerFunc),
		prefix:     r.prefix + prefix,
	}
}

func (r *Router) Add(method string, path string, handler HandlerFunc, middleware ...Middleware) {
	if _, ok := r.trees[method]; !ok {
		r.trees[method] = route.NewTree()
	}
	finalHandler := handler
	if len(middleware) > 0 {
		for i := len(middleware) - 1; i >= 0; i-- {
			finalHandler = middleware[i](finalHandler)
		}
	}
	for i := len(r.middleware) - 1; i >= 0; i-- {
		finalHandler = r.middleware[i](finalHandler)
	}
	node := r.trees[method].AddNode(r.prefix + path)
	r.handlers[node] = finalHandler
}

func (r *Router) Get(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodGet, path, handler, middleware...)
}

func (r *Router) Post(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodPost, path, handler, middleware...)
}

func (r *Router) Put(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodPut, path, handler, middleware...)
}

func (r *Router) Delete(path string, handler HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodDelete, path, handler, middleware...)
}

func (r *Router) Handle(w http.ResponseWriter, req *http.Request) (handler HandlerFunc, ctx *Context) {
	method, path := req.Method, req.URL.Path

	if _, ok := r.trees[method]; !ok {
		return
	}
	node := r.trees[method].GetNode(path)
	if node == nil {
		return
	}
	handler = r.handlers[node]
	if handler == nil {
		return
	}
	ctx = &Context{
		Writer:  w,
		Request: req,
		Params:  ParamsExtract(node.GetFullPath(), path),
	}
	return
}

func ParamsExtract(pattern string, path string) map[string]string {
	params := make(map[string]string)

	// 1. 按 "/" 分割成切片
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// 2. 长度不一致说明完全不匹配
	if len(patternParts) != len(pathParts) {
		return params
	}

	// 3. 遍历匹配
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			// 如果是 : 开头的，提取到 map
			key := strings.TrimPrefix(part, ":")
			params[key] = pathParts[i]
		} else if part != pathParts[i] {
			// 如果不是参数位且不相等，说明路径不匹配，清空并返回
			return make(map[string]string)
		}
	}

	return params
}
