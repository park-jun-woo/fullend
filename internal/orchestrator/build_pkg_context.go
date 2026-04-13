//ff:func feature=orchestrator type=util control=sequence
//ff:what buildPkgContext — pkg/generate 호출용 Fullstack + Ground 구축
package orchestrator

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/ground"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// buildPkgContext parses all SSOTs under specsDir once and derives a Ground.
// Phase008 어댑터: orchestrator 의 ParsedSSOTs 흐름과 병행 존재.
// 장기: ParseAll → Fullstack 단일화 (별도 Phase).
func buildPkgContext(specsDir string) (*fullend.Fullstack, *rule.Ground) {
	detected, _ := fullend.DetectSSOTs(specsDir)
	fs := fullend.ParseAll(specsDir, detected, nil)
	g := ground.Build(fs)
	return fs, g
}
