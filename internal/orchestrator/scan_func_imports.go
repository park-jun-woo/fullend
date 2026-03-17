//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what scanFuncImports scans SSaC files for import statements that reference func packages.

package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
)

// scanFuncImports scans SSaC files for import statements that reference func packages.
// Returns a map of package name to full import path.
func scanFuncImports(specsDir, modulePath string) (map[string]string, error) {
	result := make(map[string]string)

	ssacFiles, _ := filepath.Glob(filepath.Join(specsDir, "service", "**", "*.ssac"))
	if len(ssacFiles) == 0 {
		ssacFiles, _ = filepath.Glob(filepath.Join(specsDir, "service", "*.ssac"))
	}

	for _, f := range ssacFiles {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "import ") {
				continue
			}
			// Extract quoted path.
			q1 := strings.Index(line, "\"")
			q2 := strings.LastIndex(line, "\"")
			if q1 < 0 || q2 <= q1 {
				continue
			}
			importPath := line[q1+1 : q2]

			// Skip fullend built-in packages.
			if strings.HasPrefix(importPath, "github.com/park-jun-woo/fullend/") {
				continue
			}

			// Only consider imports within the project module.
			if !strings.HasPrefix(importPath, modulePath+"/") {
				continue
			}

			// Extract package name (last segment).
			pkg := filepath.Base(importPath)
			result[pkg] = importPath
		}
	}

	return result, nil
}
