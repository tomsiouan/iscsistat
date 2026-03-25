package main

import (
	"context"

	"iscsistat/internal/collector"
)

func (s *Server) StartCollector(ctx context.Context) {
	cfg := collector.CollectorConfig{
		CollectInterval: s.cfg.Metrics.CollectInterval,
		ISCSI:           s.cfg.ISCSI,
	}

	c := collector.New(cfg, s.logger, updateVolumeMetrics)
	c.Start(ctx)
}
