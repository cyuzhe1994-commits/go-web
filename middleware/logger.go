package middleware

import (
	"log"
	"time"

	"github.com/cyuzhe1994-commits/go-web"
)

func Logger(next go_web.HandlerFunc) go_web.HandlerFunc {
	return func(c *go_web.Context) {
		// Start timer
		t := time.Now()
		// Process request
		next(c)
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
