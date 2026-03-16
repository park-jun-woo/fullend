//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what SSaC @state가 참조하는 다이어그램 존재 및 유효한 전이 이벤트 검증
package crosscheck

import (
	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// checkGuardStates validates that SSaC guard states reference existing diagrams with valid transitions.
func checkGuardStates(funcs []ssacparser.ServiceFunc, diagramByID map[string]*statemachine.StateDiagram) []CrossError {
	var errs []CrossError
	for _, fn := range funcs {
		errs = append(errs, checkFuncGuardStates(fn, diagramByID)...)
	}
	return errs
}
