package middleware

import (
	"github.com/cyuzhe1994-commits/go-web"
)

func Cors(next go_web.HandlerFunc) go_web.HandlerFunc {
	return func(c *go_web.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}
		next(c)
	}
}
