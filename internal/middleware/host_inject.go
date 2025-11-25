package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InjectHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		hostInfo := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
		ctx := context.WithValue(c.Request.Context(), "hostInfo", hostInfo)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
