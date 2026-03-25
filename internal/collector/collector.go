package collector

import (
	"context"
	"time"

	"iscsistat/internal/config"

	"go.uber.org/zap"
)

// Config holds the collector configuration.
type CollectorConfig struct {
	CollectInterval time.Duration
	ISCSI           config.ISCSIConfig
}

// UpdateMetricsFn is a callback used to report collected metrics for a target/device pair.
type UpdateMetricsFn func(targetName, deviceName string, used, total float64)

// Collector periodically gathers iSCSI disk usage metrics and reports them via UpdateMetricsFn.
type Collector struct {
	cfg           CollectorConfig
	logger        *zap.Logger
	updateMetrics UpdateMetricsFn
}

// New creates and returns a new Collector.
func New(cfg CollectorConfig, logger *zap.Logger, updateFn UpdateMetricsFn) *Collector {
	return &Collector{
		cfg:           cfg,
		logger:        logger,
		updateMetrics: updateFn,
	}
}

// Start launches the collection loop in a background goroutine.
// It stops when the provided context is cancelled.
func (c *Collector) Start(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.CollectInterval)

	go func() {
		defer ticker.Stop()
		defer c.logger.Info("collector stopped")

		c.logger.Info("collector started", zap.Duration("interval", c.cfg.CollectInterval))

		for {
			select {
			case <-ticker.C:
				c.collect()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// collect discovers active iSCSI sessions and updates metrics for each device.
func (c *Collector) collect() {
	c.logger.Debug("collecting iSCSI metrics...")

	sessions, err := discoverISCSISessions(c.cfg.ISCSI, c.logger)
	if err != nil {
		c.logger.Error("failed to discover iSCSI sessions", zap.Error(err))
		return
	}

	if len(sessions) == 0 {
		c.logger.Debug("no active iSCSI sessions found")
		return
	}

	for _, session := range sessions {
		mountpoint, err := findMountpoint(session.DeviceName)
		if err != nil {
			c.logger.Error("failed to find mountpoint", zap.String("device", session.DeviceName), zap.Error(err))
			continue
		}

		total, used, err := getDiskStats(mountpoint, c.logger)
		if err != nil {
			c.logger.Error("failed to get disk stats", zap.String("mountpoint", mountpoint), zap.Error(err))
			continue
		}

		c.updateMetrics(session.TargetName, session.DeviceName, float64(used), float64(total))

		c.logger.Debug("collected stats",
			zap.String("target", session.TargetName),
			zap.String("device", session.DeviceName),
			zap.String("mountpoint", mountpoint),
			zap.Uint64("total", total),
			zap.Uint64("used", used),
		)
	}

	c.logger.Debug("finished collecting iSCSI metrics")
}
