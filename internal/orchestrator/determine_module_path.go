//ff:func feature=orchestrator type=util control=sequence
//ff:what determineModulePath resolves the Go module path from config, go.mod, or directory name.

package orchestrator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/projectconfig"
)

func determineModulePath(specsDir, artifactsDir string, cfg *projectconfig.ProjectConfig) string {
	// 1. Try pre-parsed config first.
	if cfg != nil && cfg.Backend.Module != "" {
		return cfg.Backend.Module
	}

	// 2. Fallback: check existing backend/go.mod.
	goModPath := filepath.Join(artifactsDir, "backend", "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "module "))
			}
		}
	}

	// 3. Last resort: derive from directory name.
	base := filepath.Base(artifactsDir)
	return base + "/backend"
}
