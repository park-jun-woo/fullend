//ff:func feature=pkg-authz type=loader control=iteration dimension=1
//ff:what loadPolicyFromPath — 디렉토리면 .rego glob concat, 파일이면 ReadFile

package authz

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// loadPolicyFromPath reads a single .rego file OR concatenates all *.rego files in a directory.
// For directory: files sorted by name for deterministic output.
func loadPolicyFromPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat OPA policy path %s: %w", path, err)
	}
	if !info.IsDir() {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read OPA policy %s: %w", path, err)
		}
		return string(data), nil
	}
	matches, err := filepath.Glob(filepath.Join(path, "*.rego"))
	if err != nil {
		return "", fmt.Errorf("glob .rego under %s: %w", path, err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no .rego files under %s", path)
	}
	sort.Strings(matches)
	var sb strings.Builder
	for _, m := range matches {
		data, rerr := os.ReadFile(m)
		if rerr != nil {
			return "", fmt.Errorf("read %s: %w", m, rerr)
		}
		sb.Write(data)
		sb.WriteByte('\n')
	}
	return sb.String(), nil
}
