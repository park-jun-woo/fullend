//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what Func 스펙 검증 — func spec 파일 수 + TODO 스텁 수 집계
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func validateFunc(specs []funcspec.FuncSpec) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindFunc)}
	if specs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "Func parse failed")
		return step
	}
	if len(specs) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no func spec files found"
		return step
	}

	// Count stubs.
	stubs := 0
	for _, s := range specs {
		if !s.HasBody {
			stubs++
		}
	}

	step.Status = reporter.Pass
	if stubs > 0 {
		step.Summary = fmt.Sprintf("%d funcs (%d TODO)", len(specs), stubs)
	} else {
		step.Summary = fmt.Sprintf("%d funcs", len(specs))
	}
	return step
}
