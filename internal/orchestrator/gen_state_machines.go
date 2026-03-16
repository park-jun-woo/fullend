//ff:func feature=orchestrator type=command control=sequence
//ff:what genStateMachines generates Go state machine code from Mermaid stateDiagram specs.

package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/gen/gogin"
	"github.com/geul-org/fullend/internal/genapi"
	"github.com/geul-org/fullend/internal/reporter"
)

func genStateMachines(specsDir, artifactsDir string, parsed *genapi.ParsedSSOTs) reporter.StepResult {
	step := reporter.StepResult{Name: "state-gen"}

	diagrams := parsed.StateDiagrams
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

	modulePath := determineModulePath(specsDir, artifactsDir, parsed.Config)
	if err := gogin.GenerateStateMachines(diagrams, artifactsDir, modulePath); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("state-gen error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d state machines generated", len(diagrams))
	return step
}
