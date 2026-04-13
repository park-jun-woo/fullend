//ff:func feature=gen-react type=generator control=sequence
//ff:what Generate — React + Vite 프론트엔드 생성 진입점
package react

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

// Config holds frontend generation configuration.
type Config struct {
	ArtifactsDir string
}

// STMLGenOutput holds STML generator output passed into React generator.
type STMLGenOutput struct {
	Deps    map[string]string
	Pages   []string
	PageOps map[string]string
}

// Generate creates React + Vite frontend from Fullstack + STML output.
func Generate(fs *fullend.Fullstack, cfg *Config, stmlOut *STMLGenOutput) error {
	var deps map[string]string
	var pages []string
	var pageOps map[string]string
	if stmlOut != nil {
		deps = stmlOut.Deps
		pages = stmlOut.Pages
		pageOps = stmlOut.PageOps
	}
	return generateFrontendSetup(cfg.ArtifactsDir, fs.OpenAPIDoc, deps, pages, pageOps)
}
