package collector

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// findMountpoint search into /proc/mounts the mount point of /dev/<deviceName>
func findMountpoint(deviceName string) (string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	target := "/dev/" + deviceName
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		if fields[0] == target {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("no mountpoint found for %s", deviceName)
}
