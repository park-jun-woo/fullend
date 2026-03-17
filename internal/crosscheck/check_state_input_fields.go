//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what @state Inputs 필드가 DDL 컬럼에 존재하는지 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// checkStateInputFields validates that @state input fields map to existing DDL columns.
func checkStateInputFields(funcs []ssacparser.ServiceFunc, diagramByID map[string]*statemachine.StateDiagram, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError

	for _, fn := range funcs {
		errs = append(errs, checkFuncStateInputs(fn, diagramByID, st)...)
	}

	return errs
}
