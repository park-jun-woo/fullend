//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=states
//ff:what 단일 SSaC 함수의 @state 입력 필드를 DDL 컬럼과 대조 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// checkFuncStateInputs validates @state input fields for a single function.
func checkFuncStateInputs(fn ssacparser.ServiceFunc, diagramByID map[string]*statemachine.StateDiagram, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError

	for _, seq := range fn.Sequences {
		if seq.Type != "state" || len(seq.Inputs) == 0 {
			continue
		}
		if _, ok := diagramByID[seq.DiagramID]; !ok {
			continue
		}

		statusField := extractStatusField(seq)
		tableName := diagramIDToTable(seq.DiagramID)
		colName := pascalToSnakeState(statusField)

		if !columnExistsInTable(st, tableName, colName) {
			errs = append(errs, CrossError{
				Rule:       "States ↔ DDL",
				Context:    fn.Name,
				Message:    fmt.Sprintf("state field %q maps to column %s.%s which does not exist", statusField, tableName, colName),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add column %s to table %s in DDL", colName, tableName),
			})
		}
	}

	return errs
}
