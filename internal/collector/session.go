package collector

import (
	"iscsistat/internal/config"
	"iscsistat/internal/iscsi"

	"go.uber.org/zap"
)

// SessionInfo holds the parsed target name and associated block device for an iSCSI session.
type SessionInfo struct {
	TargetName string // name extracted by the parser
	DeviceName string // e.g., "sda"
}

// discoverISCSISessions reads active iSCSI sessions from sysfs, parses their IQNs,
// and returns the associated block devices.
func discoverISCSISessions(cfg config.ISCSIConfig, logger *zap.Logger) ([]SessionInfo, error) {
	devices, err := iscsi.GetDevices(cfg.Parser, cfg.TrimPrefix, cfg.TrimSuffix)
	if err != nil {
		return nil, err
	}

	var sessions []SessionInfo
	for targetName, deviceName := range devices {
		logger.Debug("discovered iSCSI session",
			zap.String("target", targetName),
			zap.String("device", deviceName),
		)
		sessions = append(sessions, SessionInfo{
			TargetName: targetName,
			DeviceName: deviceName,
		})
	}

	return sessions, nil
}
