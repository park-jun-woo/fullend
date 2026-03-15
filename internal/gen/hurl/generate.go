//ff:func feature=gen-hurl type=generator
//ff:what Creates Hurl smoke tests from parsed SSOTs (public entry point).
package hurl

import "github.com/geul-org/fullend/internal/genapi"

// Generate creates Hurl smoke tests from parsed SSOTs.
func Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig) error {
	return generateHurlTests(parsed.OpenAPIDoc, cfg.ArtifactsDir, cfg.SpecsDir,
		parsed.StateDiagrams, parsed.ServiceFuncs, parsed.Policies)
}
