package main

import (
	"context"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (s *Server) StartCollector(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.Metrics.CollectInterval)

	go func() {
		defer ticker.Stop()
		defer s.logger.Info("collector background task stopped")

		s.logger.Info("collector background task started", zap.Duration("interval", s.cfg.Metrics.CollectInterval))

		for {
			select {
			case <-ticker.C:
				s.collectMetrics()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Server) collectMetrics() {
	// TODO: Dynamiser cette liste
	mounts := map[string]string{
		"n8ndata": "/var/lib/nomad/client/csi/node/synology/per-alloc/v-123/mount",
	}

	for name, path := range mounts {
		var stat syscall.Statfs_t
		err := syscall.Statfs(path, &stat)
		if err != nil {
			s.logger.Error("failed to get stats", zap.String("path", path), zap.Error(err))
			continue
		}

		all := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bfree * uint64(stat.Bsize)
		used := all - free

		volumeTotal.WithLabelValues(name).Set(float64(all))
		volumeUsage.WithLabelValues(name).Set(float64(used))
	}
}
