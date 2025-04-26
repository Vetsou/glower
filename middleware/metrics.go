package middleware

import (
	"glower/metrics"

	"github.com/gin-gonic/gin"
)

func CreateMetrics(c *gin.Context) {
	metrics.HTTPTotalRequest.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
	c.Next()
}
