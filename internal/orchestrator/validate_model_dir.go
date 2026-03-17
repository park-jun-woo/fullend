//ff:func feature=orchestrator type=rule control=sequence
//ff:what model 디렉토리 검증 — *.go 파일 존재 여부 확인
package orchestrator

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/reporter"
)

func validateModel(modelDir string) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindModel)}
	matches, _ := filepath.Glob(filepath.Join(modelDir, "*.go"))
	if len(matches) == 0 {
		step.Status = reporter.Fail
		step.Summary = "no model files found"
		return step
	}
	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d files", len(matches))
	return step
}
