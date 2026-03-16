//ff:func feature=orchestrator type=util control=selection
//ff:what returns a StepResult for a missing SSOT kind

package orchestrator

import "github.com/geul-org/fullend/internal/reporter"

// missingSSOTStep returns a StepResult for a missing SSOT kind.
func missingSSOTStep(kind SSOTKind) reporter.StepResult {
	switch kind {
	case KindFunc:
		return reporter.StepResult{Name: string(kind), Status: reporter.Skip, Summary: "no func/ directory"}
	case KindStates:
		return reporter.StepResult{Name: string(kind), Status: reporter.Skip, Summary: "no states/ directory"}
	case KindPolicy:
		return reporter.StepResult{Name: string(kind), Status: reporter.Skip, Summary: "no policy/ directory"}
	case KindScenario:
		return reporter.StepResult{
			Name: string(kind), Status: reporter.Pass, Summary: "no scenario tests",
			Errors: []string{"[WARN] tests/scenario-*.hurl 파일이 없습니다 — 시나리오 테스트를 작성하세요 (--skip scenario로 억제 가능)"},
		}
	default:
		return reporter.StepResult{Name: string(kind), Status: reporter.Fail, Summary: "required but not found"}
	}
}
