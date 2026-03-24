package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// volumeUsage defines a Prometheus Gauge metric that tracks the used space
	// for each iSCSI volume. It uses "volume" (mount path) and "device" (system device)
	// as labels for fine-grained monitoring.
	volumeUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iscsi_volume_used_bytes",
			Help: "Space used by the iSCSI volume in bytes.",
		},
		[]string{"volume", "device"},
	)
	// volumeTotal defines a Prometheus Gauge metric that tracks the total capacity
	// of each iSCSI volume. This allows calculating the percentage of usage.
	volumeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iscsi_volume_total_bytes",
			Help: "Total iSCSI volume capacity in bytes.",
		},
		[]string{"volume", "device"},
	)
)

// initMetrics registers the Prometheus metrics collectors with the default registry.
// This function must be called once during the application startup, typically in main.
func initMetrics() {
	prometheus.MustRegister(volumeUsage)
	prometheus.MustRegister(volumeTotal)
}

// updateVolumeMetrics updates the gauges for used and total volume space.
// It uses the volume name (target) and device name (e.g., sda) as label values
// to ensure granular tracking of each iSCSI session.
func updateVolumeMetrics(volume string, device string, used float64, total float64) {
	volumeUsage.WithLabelValues(volume, device).Set(used)
	volumeTotal.WithLabelValues(volume, device).Set(total)
}
