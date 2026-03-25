package parsers

import (
	"strings"
)

type DemocraticCSIParser struct{}

func (p DemocraticCSIParser) Name() string { return "democratic-csi" }

// Parse extracts the part after ":csi." in a democratic-csi IQN.
// Returns the short name and true on success, or "", false if not found.
func (p DemocraticCSIParser) Parse(iqn string) (string, bool) {
	parts := strings.SplitN(iqn, ":csi.", 2)
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}
