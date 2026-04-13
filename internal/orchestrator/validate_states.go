//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what Mermaid stateDiagram 검증 — pkg/validate/statemachine 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/internal/statemachine"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	pkgstates "github.com/park-jun-woo/fullend/pkg/validate/statemachine"
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

	detected, _ := fullend.DetectSSOTs(".")
	fs := fullend.ParseAll(".", detected, nil)
	verrs := pkgstates.Validate(fs.StatesDiags)
	for _, ve := range verrs {
		step.Errors = append(step.Errors, ve.Message)
	}

	totalTransitions := countTransitions(diagrams)
	if len(step.Errors) > 0 {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d diagrams, %d transitions", len(diagrams), totalTransitions)
	return step
}
