//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 상태 전이가 있지만 @state 시퀀스가 없는 함수 경고
package crosscheck

import (
	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// checkMissingGuards warns when a function has state transitions but no @state sequence.
func checkMissingGuards(diagrams []*statemachine.StateDiagram, funcs []ssacparser.ServiceFunc, funcNames map[string]bool) []CrossError {
	var errs []CrossError

	guardStateFuncs := collectGuardStateFuncs(funcs)

	for _, d := range diagrams {
		errs = append(errs, checkDiagramMissingGuards(d, funcNames, guardStateFuncs)...)
	}

	return errs
}
