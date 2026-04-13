//ff:func feature=orchestrator type=command control=sequence
//ff:what genStateMachines generates Go state machine code from Mermaid stateDiagram specs (pkg 경로).

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/generate/gogin"
)

func genStateMachines(specsDir, artifactsDir string, fs *fullend.Fullstack) reporter.StepResult {
	step := reporter.StepResult{Name: "state-gen"}

	diagrams := fs.StateDiagrams
	if diagrams == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "States parse failed")
		return step
	}
	if len(diagrams) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no state diagrams"
		return step
	}

	modulePath := determinePkgModulePath(specsDir, artifactsDir, fs.Manifest)
	if err := gogin.GenerateStateMachines(diagrams, artifactsDir, modulePath); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("state-gen error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d state machines generated", len(diagrams))
	return step
}
