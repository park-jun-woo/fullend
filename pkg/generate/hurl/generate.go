//ff:func feature=gen-hurl type=generator control=sequence
//ff:what Generate — Fullstack + Ground 에서 Hurl smoke 테스트 생성
package hurl

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Config holds hurl generation configuration.
type Config struct {
	ArtifactsDir string
	SpecsDir     string
}

// Generate creates Hurl smoke tests from Fullstack.
// ground is accepted for signature consistency; current implementation does not read it directly.
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
	return generateHurlTests(
		fs.OpenAPIDoc,
		cfg.ArtifactsDir,
		cfg.SpecsDir,
		fs.StateDiagrams,
		fs.ServiceFuncs,
		fs.ParsedPolicies,
	)
}
