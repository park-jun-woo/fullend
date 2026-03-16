//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what restores preserved function bodies and reports contract change warnings

package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/reporter"
)

// restorePreserved restores preserved function bodies and returns a step result.
// Returns nil if no preserves exist.
func restorePreserved(snap *contract.PreserveSnapshot) *reporter.StepResult {
	if snap == nil {
		return nil
	}
	if len(snap.FilePreserves) == 0 && len(snap.FuncPreserves) == 0 {
		return nil
	}
	warnings := contract.RestorePreserved(snap)
	preserveCount := len(snap.FilePreserves)
	for _, funcs := range snap.FuncPreserves {
		preserveCount += len(funcs)
	}
	step := reporter.StepResult{
		Name:    "preserve",
		Status:  reporter.Pass,
		Summary: fmt.Sprintf("%d preserved", preserveCount),
	}
	for _, w := range warnings {
		step.Errors = append(step.Errors,
			fmt.Sprintf("[WARN] contract changed: %s:%s (old=%s new=%s)",
				w.File, w.Function, w.OldContract, w.NewContract))
	}
	if len(warnings) > 0 {
		step.Suggestions = append(step.Suggestions, ".new 파일에서 새 계약 코드를 확인하세요")
	}
	return &step
}
