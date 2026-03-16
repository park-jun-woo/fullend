//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=openapi-ddl
//ff:what x-filter 컬럼이 DDL 테이블에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func checkXFilter(op *openapi3.Operation, st *ssacvalidator.SymbolTable, ctx string) []CrossError {
	var errs []CrossError

	raw, ok := op.Extensions["x-filter"]
	if !ok {
		return errs
	}

	var filterExt struct {
		Allowed []string `json:"allowed"`
	}
	if err := unmarshalExt(raw, &filterExt); err != nil {
		return errs
	}

	for _, col := range filterExt.Allowed {
		snake := pascalToSnake(col)
		if !columnExistsInAnyTable(snake, st) {
			table := inferTableFromCtx(op, st)
			errs = append(errs, CrossError{
				Rule:       "x-filter ↔ DDL",
				Context:    ctx,
				Message:    fmt.Sprintf("x-filter column %q (→ %s) not found in any DDL table", col, snake),
				Suggestion: fmt.Sprintf("DDL에 추가: ALTER TABLE %s ADD COLUMN %s -- TODO: 타입 지정;", table, snake),
			})
		}
	}

	return errs
}
