//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 다이어그램에서 @state가 없는 전이 이벤트 경고
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/statemachine"
)

// checkDiagramMissingGuards checks a single diagram for events without @state guards.
func checkDiagramMissingGuards(d *statemachine.StateDiagram, funcNames, guardStateFuncs map[string]bool) []CrossError {
	var errs []CrossError
	for _, event := range d.Events() {
		if funcNames[event] && !guardStateFuncs[event] {
			errs = append(errs, CrossError{
				Rule:       "States ↔ SSaC",
				Context:    event,
				Message:    fmt.Sprintf("function %q has a state transition in %s but no @state sequence", event, d.ID),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("Add @state %s sequence to %s", d.ID, event),
			})
		}
	}
	return errs
}
