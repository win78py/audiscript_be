package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)


func TimeoutLoggerMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        elapsed := time.Since(start)
        if elapsed > timeout {
            log.Printf("WARNING: Request %s %s took %v (timeout: %v)", c.Request.Method, c.Request.URL.Path, elapsed, timeout)
        }
    }
}