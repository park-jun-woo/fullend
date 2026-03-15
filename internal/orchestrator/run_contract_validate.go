//ff:func feature=orchestrator type=command
//ff:what Contract 검증 — artifacts 디렉토리의 gen/preserve 계약 검사
package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/reporter"
)

func runContractValidate(specsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "Contract"}

	// Infer artifacts dir: ../artifacts/<basename(specsDir)>
	base := filepath.Base(specsDir)
	artifactsDir := filepath.Join(filepath.Dir(specsDir), "artifacts", base)
	backendDir := filepath.Join(artifactsDir, "backend")

	if _, err := os.Stat(backendDir); os.IsNotExist(err) {
		step.Status = reporter.Skip
		step.Summary = "no artifacts"
		return step
	}

	funcs, err := contract.ScanDir(artifactsDir)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, err.Error())
		return step
	}

	if len(funcs) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no directives"
		return step
	}

	funcs = contract.Verify(specsDir, funcs)
	gen, preserve, broken, orphan := contract.Summary(funcs)

	parts := []string{}
	if gen > 0 {
		parts = append(parts, fmt.Sprintf("%d gen", gen))
	}
	if preserve > 0 {
		parts = append(parts, fmt.Sprintf("%d preserve", preserve))
	}
	if broken > 0 {
		parts = append(parts, fmt.Sprintf("%d broken", broken))
	}
	if orphan > 0 {
		parts = append(parts, fmt.Sprintf("%d orphan", orphan))
	}
	step.Summary = strings.Join(parts, ", ")

	if broken > 0 || orphan > 0 {
		step.Status = reporter.Fail
		for _, f := range funcs {
			if f.Status == "broken" || f.Status == "orphan" {
				step.Errors = append(step.Errors, fmt.Sprintf("%s: %s %s — %s", f.Status, f.File, f.Function, f.Detail))
			}
		}
	} else {
		step.Status = reporter.Pass
	}

	return step
}
