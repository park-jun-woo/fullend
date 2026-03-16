//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what stmlMatchAttr returns the STML attribute that references the operationId.

package orchestrator

import (
	"bufio"
	"os"
	"strings"
)

// stmlMatchAttr returns the STML attribute that references the operationId (e.g. "data-fetch", "data-action").
func stmlMatchAttr(filePath, opID string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return "data-fetch"
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, opID) {
			continue
		}
		if strings.Contains(line, "data-action") {
			return "data-action"
		}
		return "data-fetch"
	}
	return "data-fetch"
}
