package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	volumeUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iscsi_volume_used_bytes",
			Help: "Space used by the iSCSI volume in bytes.",
		},
		[]string{"volume_name"},
	)
	volumeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iscsi_volume_total_bytes",
			Help: "Total iSCSI volume capacity in bytes.",
		},
		[]string{"volume_name"},
	)
)

func init() {
	prometheus.MustRegister(volumeUsage)
	prometheus.MustRegister(volumeTotal)
}

func (s *Server) handleMetrics() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
