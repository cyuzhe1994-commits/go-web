package public

type HandlerFunc func(*Context)

type Middleware func(HandlerFunc) HandlerFunc
