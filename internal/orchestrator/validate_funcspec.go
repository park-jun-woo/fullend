//ff:func feature=orchestrator type=rule control=sequence dimension=1
//ff:what Func 스펙 검증 — func spec 파일 수 + TODO 스텁 수 집계
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func validateFunc(projectSpecs, fullendSpecs []funcspec.FuncSpec) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindFunc)}
	if projectSpecs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "Func parse failed")
		return step
	}

	if errs := checkBuiltinOverride(projectSpecs, fullendSpecs); len(errs) > 0 {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, errs...)
		return step
	}

	if len(projectSpecs) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no func spec files found"
		return step
	}

	stubs := countStubs(projectSpecs)
	step.Status = reporter.Pass
	if stubs > 0 {
		step.Summary = fmt.Sprintf("%d funcs (%d TODO)", len(projectSpecs), stubs)
	} else {
		step.Summary = fmt.Sprintf("%d funcs", len(projectSpecs))
	}
	return step
}
