//ff:func feature=orchestrator type=rule control=iteration dimension=1
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

	// Check built-in package override.
	if len(fullendSpecs) > 0 {
		builtinPkgs := make(map[string]bool)
		for _, s := range fullendSpecs {
			builtinPkgs[s.Package] = true
		}
		for _, s := range projectSpecs {
			if builtinPkgs[s.Package] {
				step.Status = reporter.Fail
				step.Errors = append(step.Errors, fmt.Sprintf(
					"func/%s: built-in 패키지 %q를 override할 수 없습니다. 커스텀 패키지명을 사용하세요 (예: func/my%s/)",
					s.Package, s.Package, s.Package))
			}
		}
		if step.Status == reporter.Fail {
			return step
		}
	}

	if len(projectSpecs) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no func spec files found"
		return step
	}

	// Count stubs.
	stubs := 0
	for _, s := range projectSpecs {
		if !s.HasBody {
			stubs++
		}
	}

	step.Status = reporter.Pass
	if stubs > 0 {
		step.Summary = fmt.Sprintf("%d funcs (%d TODO)", len(projectSpecs), stubs)
	} else {
		step.Summary = fmt.Sprintf("%d funcs", len(projectSpecs))
	}
	return step
}
