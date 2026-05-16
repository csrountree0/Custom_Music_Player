package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"musicapp/backend/pkg/apierror"
)

var (
	limiters sync.Map
)

func getLimiter(ip string) *rate.Limiter {
	if v, ok := limiters.Load(ip); ok {
		return v.(*rate.Limiter)
	}
	l := rate.NewLimiter(rate.Limit(1), 60)
	limiters.Store(ip, l)
	return l
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !getLimiter(ip).Allow() {
			apierror.New(429, "RATE_LIMITED", "too many requests").Respond(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
