package private

import (
	"glower/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetrics(e *gin.Engine) {
	reg := metrics.CreateRegistry()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	metricsGroup := e.Group("/metrics")
	metricsGroup.GET("/", gin.WrapH(promHandler))
}
