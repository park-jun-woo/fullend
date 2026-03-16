//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 상태 다이어그램을 SSaC, DDL, OpenAPI와 교차 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckStates validates state diagrams against SSaC, DDL, and OpenAPI.
func CheckStates(diagrams []*statemachine.StateDiagram, funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, doc *openapi3.T) []CrossError {
	var errs []CrossError

	if len(diagrams) == 0 {
		return errs
	}

	diagramByID := make(map[string]*statemachine.StateDiagram)
	for _, d := range diagrams {
		diagramByID[d.ID] = d
	}

	funcNames := make(map[string]bool)
	for _, fn := range funcs {
		funcNames[fn.Name] = true
	}

	errs = append(errs, checkTransitionEvents(diagrams, funcNames)...)
	errs = append(errs, checkGuardStates(funcs, diagramByID)...)
	errs = append(errs, checkMissingGuards(diagrams, funcs, funcNames)...)

	if st != nil {
		errs = append(errs, checkStateInputFields(funcs, diagramByID, st)...)
	}

	return errs
}
