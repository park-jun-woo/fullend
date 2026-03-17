//ff:func feature=gen-react type=generator control=sequence
//ff:what React + Vite 프론트엔드를 SSOT에서 생성한다

package react

import (
	"github.com/park-jun-woo/fullend/internal/genapi"
)

// Generate creates React + Vite frontend from parsed SSOTs.
func Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig, stmlOut *genapi.STMLGenOutput) error {
	var deps map[string]string
	var pages []string
	var pageOps map[string]string
	if stmlOut != nil {
		deps = stmlOut.Deps
		pages = stmlOut.Pages
		pageOps = stmlOut.PageOps
	}
	return generateFrontendSetup(cfg.ArtifactsDir, parsed.OpenAPIDoc, deps, pages, pageOps)
}
