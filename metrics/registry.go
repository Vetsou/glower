package metrics

import "github.com/prometheus/client_golang/prometheus"

func CreateRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(HTTPTotalRequest)

	return reg
}
