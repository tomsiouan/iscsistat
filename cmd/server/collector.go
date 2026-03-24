package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Regex to extract size and used from df output like "  1073741824  536870912"
var dfStatsRegex = regexp.MustCompile(`^\s*(\d+)\s+(\d+)`)

// Regex to extract Target Name from iscsiadm output line like "Target: iqn.2000-01.com.synology:csi.hello (non-flash)"
var targetNameRegex = regexp.MustCompile(`Target:\s+\S+:csi\.([^\s(]+)`)

// Regex to extract SCSI disk from iscsiadm output line like "Attached scsi disk sdX"
var attachedDiskRegex = regexp.MustCompile(`Attached scsi disk (sd[a-z]+)`)

type SessionInfo struct {
	TargetName string // e.g., "volume_name"
	DeviceName string // e.g., "sda"
}

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
	s.logger.Debug("collecting iSCSI metrics...")

	sessions, err := discoverISCSISessions(s.logger)
	if err != nil {
		s.logger.Error("failed to discover iSCSI sessions", zap.Error(err))
		return
	}

	if len(sessions) == 0 {
		s.logger.Debug("no active iSCSI sessions found.")
		return
	}

	for _, session := range sessions {
		devicePath := "/dev/" + session.DeviceName
		size, used, err := getDiskStats(devicePath, s.logger)
		if err != nil {
			s.logger.Error("failed to get disk stats", zap.String("device", devicePath), zap.Error(err))
			continue
		}

		updateVolumeMetrics(session.TargetName, session.DeviceName, float64(used), float64(size))

		s.logger.Debug("collected iSCSI volume stats",
			zap.String("target", session.TargetName),
			zap.String("device", session.DeviceName),
			zap.Uint64("total_bytes", size),
			zap.Uint64("used_bytes", used),
		)
	}
	s.logger.Debug("finished collecting iSCSI metrics.")
}

// discoverISCSISessions uses `iscsiadm` to find active iSCSI sessions and their attached SCSI devices.
func discoverISCSISessions(logger *zap.Logger) ([]SessionInfo, error) {
	var sessions []SessionInfo
	foundDisks := make(map[string]map[string]struct{})

	cmd := exec.Command("iscsiadm", "-m", "session", "-P", "3")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			logger.Error("iscsiadm stderr during discovery", zap.ByteString("stderr", stderr.Bytes()))
		}
		return nil, fmt.Errorf("failed to run iscsiadm command: %w", err)
	}

	var currentTargetName string
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if matches := targetNameRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentTargetName = matches[1]
			if _, ok := foundDisks[currentTargetName]; !ok {
				foundDisks[currentTargetName] = make(map[string]struct{})
			}
			continue
		}

		if currentTargetName != "" {
			if matches := attachedDiskRegex.FindStringSubmatch(line); len(matches) > 1 {
				deviceName := matches[1]
				foundDisks[currentTargetName][deviceName] = struct{}{}
			}
		}
	}

	for targetName, disks := range foundDisks {
		for deviceName := range disks {
			sessions = append(sessions, SessionInfo{
				TargetName: targetName,
				DeviceName: deviceName,
			})
		}
	}

	return sessions, nil
}

// getDiskStats executes `df` for a given device path and returns total and used bytes.
func getDiskStats(devicePath string, logger *zap.Logger) (totalBytes, usedBytes uint64, err error) {
	cmd := exec.Command("df", "-B1", "--output=size,used", devicePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			logger.Error("df stderr during stats collection", zap.String("device", devicePath), zap.ByteString("stderr", stderr.Bytes()))
		}
		return 0, 0, fmt.Errorf("failed to run df command for %s: %w", devicePath, err)
	}

	dfOutput := stdout.String()
	lines := strings.Split(dfOutput, "\n")
	var statsLine string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "1B-blocks") && !strings.HasPrefix(trimmedLine, "Size") {
			statsLine = trimmedLine
			break
		}
	}

	if statsLine == "" {
		return 0, 0, fmt.Errorf("df command for %s returned no usable stats: %s", devicePath, dfOutput)
	}

	if matches := dfStatsRegex.FindStringSubmatch(statsLine); len(matches) > 2 {
		sizeStr := matches[1]
		usedStr := matches[2]

		size, parseErr := strconv.ParseUint(sizeStr, 10, 64)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("failed to parse size '%s' from df output for %s: %w", sizeStr, devicePath, parseErr)
		}
		used, parseErr := strconv.ParseUint(usedStr, 10, 64)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("failed to parse used '%s' from df output for %s: %w", usedStr, devicePath, parseErr)
		}

		return size, used, nil
	}

	return 0, 0, fmt.Errorf("failed to parse stats line '%s' from df output for %s", statsLine, devicePath)
}
