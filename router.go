package go_web

import (
	"net/http"
	"strings"

	"github.com/cyuzhe1994-commits/go-web/public"
	"github.com/cyuzhe1994-commits/go-web/route"
)

type Router struct {
	trees      map[string]*route.Tree // 在新建 group 后，树结构是共享的
	middleware []public.Middleware    // 在新建 group 后，中间件是复制的，保证每个 group 可以独立添加中间件
	prefix     string
}

func NewRouter() *Router {
	return &Router{
		trees:      make(map[string]*route.Tree),
		middleware: make([]public.Middleware, 0),
		prefix:     "",
	}
}

func (r *Router) Use(middleware public.Middleware) {
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) Group(prefix string) *Router {

	newMiddleware := make([]public.Middleware, len(r.middleware))
	copy(newMiddleware, r.middleware)

	return &Router{
		trees:      r.trees,
		middleware: newMiddleware,
		prefix:     r.prefix + prefix,
	}
}

func (r *Router) Add(method string, path string, handler public.HandlerFunc, middleware ...public.Middleware) {
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
	r.trees[method].AddNode(r.prefix+path, finalHandler)
}

func (r *Router) Get(path string, handler public.HandlerFunc, middleware ...public.Middleware) {
	r.Add(http.MethodGet, path, handler, middleware...)
}

func (r *Router) Post(path string, handler public.HandlerFunc, middleware ...public.Middleware) {
	r.Add(http.MethodPost, path, handler, middleware...)
}

func (r *Router) Put(path string, handler public.HandlerFunc, middleware ...public.Middleware) {
	r.Add(http.MethodPut, path, handler, middleware...)
}

func (r *Router) Delete(path string, handler public.HandlerFunc, middleware ...public.Middleware) {
	r.Add(http.MethodDelete, path, handler, middleware...)
}

func (r *Router) Handle(w http.ResponseWriter, req *http.Request) {
	method, path := req.Method, req.URL.Path
	if _, ok := r.trees[method]; !ok {
		http.NotFound(w, req)
		return
	}
	node := r.trees[method].GetNode(path)
	if node == nil {
		http.NotFound(w, req)
		return
	}
	handler := node.GetHandler()
	if handler == nil {
		http.NotFound(w, req)
		return
	}
	ctx := &public.Context{
		Writer:  w,
		Request: req,
		Params:  ParamsExtract(node.GetFullPath(), path),
	}
	handler(ctx)
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
