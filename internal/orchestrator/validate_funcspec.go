//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what Func 스펙 검증 — pkg/validate/funcspec 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	pkgfunc "github.com/park-jun-woo/fullend/pkg/validate/funcspec"
)

func validateFunc(projectSpecs, fullendSpecs []funcspec.FuncSpec) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindFunc)}
	if projectSpecs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "Func parse failed")
		return step
	}
	if len(projectSpecs) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no func spec files found"
		return step
	}

	// Use pkg/validate/funcspec for F-1 check via pkg types
	detected, _ := fullend.DetectSSOTs(".")
	fs := fullend.ParseAll(".", detected, nil)
	verrs := pkgfunc.Validate(fs.ProjectFuncSpecs, fs.FullendPkgSpecs)
	for _, ve := range verrs {
		step.Errors = append(step.Errors, fmt.Sprintf("%s: %s", ve.Rule, ve.Message))
	}

	stubs := countStubs(projectSpecs)
	if len(step.Errors) > 0 {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	if stubs > 0 {
		step.Summary = fmt.Sprintf("%d funcs (%d TODO)", len(projectSpecs), stubs)
	} else {
		step.Summary = fmt.Sprintf("%d funcs", len(projectSpecs))
	}
	return step
}
