//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what cursor 페이지네이션과 x-sort 제약 조건 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkCursorSort validates cursor pagination + x-sort constraints.
func checkCursorSort(op *openapi3.Operation, st *ssacvalidator.SymbolTable, ctx string) []CrossError {
	var errs []CrossError

	pagRaw, ok := op.Extensions["x-pagination"]
	if !ok {
		return errs
	}
	var pagExt struct {
		Style string `json:"style"`
	}
	if err := unmarshalExt(pagRaw, &pagExt); err != nil || pagExt.Style != "cursor" {
		return errs
	}

	sortRaw, ok := op.Extensions["x-sort"]
	if !ok {
		return errs
	}
	var sortExt struct {
		Allowed   []string `json:"allowed"`
		Default   string   `json:"default"`
		Direction string   `json:"direction"`
	}
	if err := unmarshalExt(sortRaw, &sortExt); err != nil {
		return errs
	}

	if len(sortExt.Allowed) > 1 {
		errs = append(errs, CrossError{
			Rule:    "x-pagination ↔ x-sort",
			Context: ctx,
			Message: fmt.Sprintf("cursor 모드에서 x-sort allowed가 %d개 — 런타임 정렬 전환은 cursor를 깨뜨립니다", len(sortExt.Allowed)),
			Level:   "ERROR",
		})
		return errs
	}

	defaultCol := sortExt.Default
	if defaultCol == "" && len(sortExt.Allowed) == 1 {
		defaultCol = sortExt.Allowed[0]
	}
	if defaultCol != "" {
		tableName := inferTableFromCtx(op, st)
		if tableName != "???" && !isUniqueColumn(defaultCol, tableName, st) {
			errs = append(errs, CrossError{
				Rule:       "x-pagination ↔ x-sort ↔ DDL",
				Context:    ctx,
				Message:    fmt.Sprintf("cursor 모드의 x-sort default %q — DDL %s에서 UNIQUE가 아닙니다. 중복값 시 cursor가 깨집니다", defaultCol, tableName),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("DDL에 UNIQUE 제약 추가: ALTER TABLE %s ADD CONSTRAINT uniq_%s_%s UNIQUE (%s);", tableName, tableName, defaultCol, defaultCol),
			})
		}
	}

	return errs
}
