//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what 전이 이벤트에 매칭되는 SSaC 함수 존재 여부 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// checkTransitionEvents validates that transition events have matching SSaC functions.
func checkTransitionEvents(diagrams []*statemachine.StateDiagram, funcNames map[string]bool) []CrossError {
	var errs []CrossError
	for _, d := range diagrams {
		errs = append(errs, checkDiagramEvents(d, funcNames)...)
	}
	return errs
}
