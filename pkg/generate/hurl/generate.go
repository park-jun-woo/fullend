//ff:func feature=gen-hurl type=generator control=sequence
//ff:what Generate — Hurl smoke 테스트 생성 진입점 (Phase004 stub)
package hurl

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Config holds hurl generation configuration.
type Config struct {
	ArtifactsDir string
	SpecsDir     string
}

// Generate creates Hurl smoke tests from Fullstack.
// STUB — Phase004 후속 작업에서 내부 함수 타입 수렴 및 Toulmin 시나리오 ordering 적용 예정.
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
	return fmt.Errorf("pkg/generate/hurl.Generate 는 아직 활성화되지 않았습니다 (Phase004 후속 작업)")
}
