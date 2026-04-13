//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what Generate — Go+Gin 백엔드 코드 생성 진입점 (Phase004 stub)

package gogin

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Config holds generation configuration for the Go+Gin backend.
type Config struct {
	ArtifactsDir string
	SpecsDir     string
	ModulePath   string
}

// Generate creates Go+Gin backend code from a parsed Fullstack.
// STUB — Phase004 의 다른 Step 에서 내부 함수 타입을 pkg 로 수렴한 뒤 실제 로직 작성.
// 현 시점에는 pkg/generate/gogin 이 빌드만 되도록 뼈대 유지.
func (g *GoGin) Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
	return fmt.Errorf("pkg/generate/gogin.Generate 는 아직 활성화되지 않았습니다 (Phase004 후속 작업)")
}
