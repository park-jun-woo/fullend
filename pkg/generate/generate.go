//ff:func feature=rule type=command control=sequence
//ff:what Generate — SSOT로부터 코드 산출물을 생성하는 최상위 엔트리포인트
package generate

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

// Generate orchestrates the full code generation pipeline.
// Returns collected errors from external tools and internal generators.
func Generate(fs *fullend.Fullstack, specsDir, artifactsDir string) []string {
	var errs []string
	errs = append(errs, runOapiCodegen(specsDir, artifactsDir)...)
	errs = append(errs, runSqlc(specsDir, artifactsDir)...)
	errs = append(errs, generateBackend(fs, artifactsDir)...)
	errs = append(errs, generateFrontend(fs, artifactsDir)...)
	errs = append(errs, generateHurl(fs, artifactsDir)...)
	return errs
}
