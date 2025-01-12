package middleware

import (
	"glower/metrics"

	"github.com/gin-gonic/gin"
)

func CountRequest(c *gin.Context) {
	metrics.HTTPTotalRequest.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
	c.Next()
}
