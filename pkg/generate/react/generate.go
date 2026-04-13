//ff:func feature=gen-react type=generator control=sequence
//ff:what Generate — React + Vite 프론트엔드 생성 진입점 (Phase004 stub)
package react

import (
	"fmt"

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
// STUB — Phase004 후속 작업에서 활성화.
func Generate(fs *fullend.Fullstack, cfg *Config, stmlOut *STMLGenOutput) error {
	return fmt.Errorf("pkg/generate/react.Generate 는 아직 활성화되지 않았습니다 (Phase004 후속 작업)")
}
