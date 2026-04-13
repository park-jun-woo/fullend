//ff:func feature=orchestrator type=util control=sequence
//ff:what determinePkgModulePath — pkg manifest 기반 Go module path 결정 (determineModulePath 의 pkg 버전)

package orchestrator

import (
	"path/filepath"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

func determinePkgModulePath(specsDir, artifactsDir string, cfg *manifest.ProjectConfig) string {
	if cfg != nil && cfg.Backend.Module != "" {
		return cfg.Backend.Module
	}
	goModPath := filepath.Join(artifactsDir, "backend", "go.mod")
	if mod := moduleFromGoMod(goModPath); mod != "" {
		return mod
	}
	base := filepath.Base(artifactsDir)
	return base + "/backend"
}
