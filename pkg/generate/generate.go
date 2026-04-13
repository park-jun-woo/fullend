//ff:func feature=genapi type=command control=sequence
//ff:what Generate — Fullstack + Ground 에서 backend + frontend + hurl 산출물 생성
package generate

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/generate/gogin"
	"github.com/park-jun-woo/fullend/pkg/generate/hurl"
	"github.com/park-jun-woo/fullend/pkg/generate/react"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Config holds top-level generation configuration.
type Config struct {
	ArtifactsDir string
	SpecsDir     string
	ModulePath   string
}

// Generate creates all artifacts from Fullstack + Ground.
// STUB — Phase004 후속 작업에서 각 하위 generator 활성화 후 전체 배선.
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config, stmlOut *react.STMLGenOutput) error {
	backend := selectBackend()
	backendCfg := &gogin.Config{
		ArtifactsDir: cfg.ArtifactsDir,
		SpecsDir:     cfg.SpecsDir,
		ModulePath:   cfg.ModulePath,
	}
	if err := backend.Generate(fs, ground, backendCfg); err != nil {
		return err
	}
	reactCfg := &react.Config{ArtifactsDir: cfg.ArtifactsDir}
	if err := react.Generate(fs, reactCfg, stmlOut); err != nil {
		return err
	}
	hurlCfg := &hurl.Config{ArtifactsDir: cfg.ArtifactsDir, SpecsDir: cfg.SpecsDir}
	return hurl.Generate(fs, ground, hurlCfg)
}
