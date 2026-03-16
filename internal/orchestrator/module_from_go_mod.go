//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what reads go.mod and extracts the module path

package orchestrator

import (
	"os"
	"strings"
)

// moduleFromGoMod reads go.mod and extracts the module path.
func moduleFromGoMod(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}
