//ff:func feature=orchestrator type=rule control=sequence
//ff:what fullend.yaml 설정 검증 — 프로젝트 메타데이터 + 백엔드/프론트엔드 구성 확인
package orchestrator

import (
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/projectconfig"
	"github.com/geul-org/fullend/internal/reporter"
)

func validateConfig(path string, cfg *projectconfig.ProjectConfig) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindConfig)}
	if cfg == nil {
		// Parse failed in ParseAll; try again for error message.
		var err error
		cfg, err = projectconfig.Load(filepath.Dir(path))
		if err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, err.Error())
			return step
		}
	}
	step.Status = reporter.Pass
	parts := []string{cfg.Metadata.Name}
	if cfg.Backend.Module != "" {
		parts = append(parts, cfg.Backend.Lang+"/"+cfg.Backend.Framework)
	}
	if cfg.Frontend.Name != "" {
		parts = append(parts, cfg.Frontend.Lang+"/"+cfg.Frontend.Framework)
	}
	step.Summary = strings.Join(parts, ", ")
	return step
}
