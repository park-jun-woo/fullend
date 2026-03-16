//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what x-sort 컬럼이 DDL 테이블에 존재하고 인덱스가 있는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func checkXSort(op *openapi3.Operation, st *ssacvalidator.SymbolTable, ctx string) []CrossError {
	var errs []CrossError

	raw, ok := op.Extensions["x-sort"]
	if !ok {
		return errs
	}

	var sortExt struct {
		Allowed []string `json:"allowed"`
	}
	if err := unmarshalExt(raw, &sortExt); err != nil {
		return errs
	}

	for _, col := range sortExt.Allowed {
		snake := pascalToSnake(col)
		if !columnExistsInAnyTable(snake, st) {
			table := inferTableFromCtx(op, st)
			errs = append(errs, CrossError{
				Rule:       "x-sort ↔ DDL",
				Context:    ctx,
				Message:    fmt.Sprintf("x-sort column %q (→ %s) not found in any DDL table", col, snake),
				Suggestion: fmt.Sprintf("DDL에 추가: ALTER TABLE %s ADD COLUMN %s -- TODO: 타입 지정;", table, snake),
			})
		} else if !columnHasUsableIndex(snake, st) {
			table := findTableWithColumn(snake, st)
			errs = append(errs, CrossError{
				Rule:       "x-sort ↔ DDL index",
				Context:    ctx,
				Message:    fmt.Sprintf("x-sort column %q (→ %s) has no index (performance)", col, snake),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("DDL에 추가: CREATE INDEX idx_%s_%s ON %s(%s);", table, snake, table, snake),
			})
		}
	}

	return errs
}
