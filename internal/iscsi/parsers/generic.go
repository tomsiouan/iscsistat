package parsers

import (
	"strings"
)

type GenericParser struct{}

func (p GenericParser) Name() string { return "generic" }

// Parse extracts the part after the last colon in the IQN.
// Returns the short name and true on success, or "", false if not found.
func (p GenericParser) Parse(iqn string) (string, bool) {
	idx := strings.LastIndex(iqn, ":")
	if idx < 0 || idx == len(iqn)-1 {
		return "", false
	}
	return iqn[idx+1:], true
}
