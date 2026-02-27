package middleware

import (
	"log"
	"time"

	"github.com/go-web/public"
)

func Logger(next public.HandlerFunc) public.HandlerFunc {
	return func(c *public.Context) {
		// Start timer
		t := time.Now()
		// Process request
		next(c)
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
