package collector

import (
	"fmt"
	"syscall"

	"go.uber.org/zap"
)

// getDiskStats uses syscall.Statfs to retrieve total and used bytes for the given device path.
func getDiskStats(mountpoint string, logger *zap.Logger) (totalBytes uint64, usedBytes uint64, err error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mountpoint, &stat); err != nil {
		return 0, 0, fmt.Errorf("statfs failed on %s: %w", mountpoint, err)
	}

	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bfree) * uint64(stat.Bsize)
	used := total - free

	logger.Debug("Statfs details",
		zap.String("mountpoint", mountpoint),
		zap.Uint64("blocks", uint64(stat.Blocks)),
		zap.Uint64("bfree", uint64(stat.Bfree)),
		zap.Uint64("bavail", uint64(stat.Bavail)),
		zap.Uint64("bsize", uint64(stat.Bsize)),
		zap.Uint64("total_bytes", total),
		zap.Uint64("used_bytes", used),
	)

	return total, used, nil
}
