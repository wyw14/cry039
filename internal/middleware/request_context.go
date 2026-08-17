package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync/atomic"
)

var next uint64

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("office-%08d", atomic.AddUint64(&next, 1))
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
