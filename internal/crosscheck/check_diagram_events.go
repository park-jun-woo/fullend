//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what 단일 다이어그램의 이벤트가 SSaC 함수에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// checkDiagramEvents validates events of a single diagram against SSaC functions.
func checkDiagramEvents(d *statemachine.StateDiagram, funcNames map[string]bool) []CrossError {
	var errs []CrossError
	for _, event := range d.Events() {
		if !funcNames[event] {
			errs = append(errs, CrossError{
				Rule:       "States ↔ SSaC",
				Context:    fmt.Sprintf("%s.%s", d.ID, event),
				Message:    fmt.Sprintf("transition event %q has no matching SSaC function", event),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add SSaC function %s or remove transition from states/%s.md", event, d.ID),
			})
		}
	}
	return errs
}
