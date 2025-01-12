package initializers

import (
	"glower/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

func CreateMetricsRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(metrics.HTTPTotalRequest)

	return reg
}
