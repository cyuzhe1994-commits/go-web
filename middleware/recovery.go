package middleware

import (
	"github.com/cyuzhe1994-commits/go-web"
)

func Recovery(next go_web.HandlerFunc) go_web.HandlerFunc {
	return func(ctx *go_web.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(500, map[string]interface{}{
					"error": err,
				})
			}
		}()
		next(ctx)
	}
}
