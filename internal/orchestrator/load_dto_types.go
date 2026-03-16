//ff:func feature=orchestrator type=loader control=iteration
//ff:what 모델 파일에서 @dto 타입 로드 — model/*.go 스캔
package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
)

// loadDTOTypes scans model/*.go files for types preceded by a // @dto comment.
func loadDTOTypes(modelDir string) map[string]bool {
	dtoTypes := make(map[string]bool)
	if modelDir == "" {
		return dtoTypes
	}
	matches, _ := filepath.Glob(filepath.Join(modelDir, "*.go"))
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		dtoNext := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "// @dto" || strings.HasPrefix(trimmed, "// @dto ") {
				dtoNext = true
				continue
			}
			if dtoNext && strings.HasPrefix(trimmed, "type ") {
				parts := strings.Fields(trimmed)
				if len(parts) >= 2 {
					dtoTypes[parts[1]] = true
				}
				dtoNext = false
			} else if dtoNext && trimmed != "" && !strings.HasPrefix(trimmed, "//") {
				dtoNext = false
			}
		}
	}
	return dtoTypes
}
