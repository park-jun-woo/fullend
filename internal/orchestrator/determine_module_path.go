//ff:func feature=orchestrator type=util control=sequence
//ff:what determineModulePath resolves the Go module path from config, go.mod, or directory name.

package orchestrator

import (
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func determineModulePath(specsDir, artifactsDir string, cfg *projectconfig.ProjectConfig) string {
	// 1. Try pre-parsed config first.
	if cfg != nil && cfg.Backend.Module != "" {
		return cfg.Backend.Module
	}

	// 2. Fallback: check existing backend/go.mod.
	goModPath := filepath.Join(artifactsDir, "backend", "go.mod")
	if mod := moduleFromGoMod(goModPath); mod != "" {
		return mod
	}

	// 3. Last resort: derive from directory name.
	base := filepath.Base(artifactsDir)
	return base + "/backend"
}
