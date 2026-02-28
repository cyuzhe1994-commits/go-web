package go_web

type HandlerFunc func(*Context)

type Middleware func(HandlerFunc) HandlerFunc
