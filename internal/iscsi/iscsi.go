package iscsi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetDevices returns a map of short target names to block device names (e.g. "mytarget" → "sda")
// by reading iSCSI session info from sysfs and parsing IQNs with the given parser.
func GetDevices(parserName, trimPrefix, trimSuffix string) (map[string]string, error) {
	parser, ok := GetParser(parserName)
	if !ok {
		return nil, fmt.Errorf("parser iSCSI inconnu: %s", parserName)
	}

	result := make(map[string]string)

	sessions, err := filepath.Glob("/sys/class/iscsi_session/session*")
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		raw, err := os.ReadFile(session + "/targetname")
		if err != nil {
			continue
		}
		iqn := strings.TrimSpace(string(raw))

		shortName, ok := parser.Parse(iqn)
		if !ok {
			continue
		}

		if trimPrefix != "" {
			shortName = strings.TrimPrefix(shortName, trimPrefix)
		}
		if trimSuffix != "" {
			shortName = strings.TrimSuffix(shortName, trimSuffix)
		}

		blocks, err := filepath.Glob(session + "/device/target*/*/block/*")
		if err != nil || len(blocks) == 0 {
			continue
		}

		result[shortName] = filepath.Base(blocks[0])
	}

	return result, nil
}
