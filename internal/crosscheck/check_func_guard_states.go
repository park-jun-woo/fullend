//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what 단일 SSaC 함수의 @state 참조 다이어그램 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// checkFuncGuardStates validates guard states for a single SSaC function.
func checkFuncGuardStates(fn ssacparser.ServiceFunc, diagramByID map[string]*statemachine.StateDiagram) []CrossError {
	var errs []CrossError
	for _, seq := range fn.Sequences {
		if seq.Type != "state" {
			continue
		}
		diagramID := seq.DiagramID
		if _, ok := diagramByID[diagramID]; !ok {
			errs = append(errs, CrossError{
				Rule:       "States ↔ SSaC",
				Context:    fn.Name,
				Message:    fmt.Sprintf("@state references diagram %q which does not exist", diagramID),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Create states/%s.md with a Mermaid stateDiagram", diagramID),
			})
			continue
		}

		d := diagramByID[diagramID]
		validStates := d.ValidFromStates(fn.Name)
		if len(validStates) == 0 {
			errs = append(errs, CrossError{
				Rule:       "States ↔ SSaC",
				Context:    fn.Name,
				Message:    fmt.Sprintf("function %q is not a valid transition event in diagram %q", fn.Name, diagramID),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add transition to states/%s.md: someState --> targetState: %s", diagramID, fn.Name),
			})
		}
	}
	return errs
}
