package middleware

import "github.com/go-web/public"

func Recovery(next public.HandlerFunc) public.HandlerFunc {
	return func(ctx *public.Context) {
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
