//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what Mermaid stateDiagram 검증 — 다이어그램 수 + 전이 수 집계
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/internal/statemachine"
)

func validateStates(diagrams []*statemachine.StateDiagram, parseErr error) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindStates)}
	if diagrams == nil {
		step.Status = reporter.Fail
		if parseErr != nil {
			step.Errors = append(step.Errors, parseErr.Error())
		} else {
			step.Errors = append(step.Errors, "States parse failed")
		}
		return step
	}
	if len(diagrams) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no state diagrams found"
		return step
	}

	totalTransitions := 0
	for _, d := range diagrams {
		totalTransitions += len(d.Transitions)
	}
	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d diagrams, %d transitions", len(diagrams), totalTransitions)
	return step
}
