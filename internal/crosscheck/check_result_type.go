//ff:func feature=crosscheck type=rule control=sequence topic=ssac-ddl
//ff:what SSaC @result 타입이 DDL 테이블에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/jinzhu/inflection"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func checkResultType(seq ssacparser.Sequence, st *ssacvalidator.SymbolTable, ctx string, seqIdx int, dtoTypes map[string]bool) []CrossError {
	var errs []CrossError

	typeName := normalizeTypeName(seq.Result.Type)

	// Skip primitive Go types.
	if primitiveTypes[typeName] {
		return errs
	}

	// Result type must be singular (e.g. "Gig" not "Gigs").
	if inflection.Singular(typeName) != typeName {
		errs = append(errs, CrossError{
			Rule:       "SSaC @result ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("seq[%d] @result type %q should be singular (expected %q)", seqIdx, seq.Result.Type, inflection.Singular(typeName)),
			Level:      "WARNING",
			Suggestion: fmt.Sprintf("@result type을 단수형 %q으로 변경", inflection.Singular(typeName)),
		})
	}

	// Skip @dto types (no DDL table).
	if dtoTypes != nil && dtoTypes[typeName] {
		return errs
	}

	tableName := modelToTable(typeName)

	if _, ok := st.DDLTables[tableName]; !ok {
		errs = append(errs, CrossError{
			Rule:       "SSaC @result ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("seq[%d] @result type %q has no matching DDL table %q", seqIdx, seq.Result.Type, tableName),
			Level:      "WARNING",
			Suggestion: fmt.Sprintf("DDL에 추가: CREATE TABLE %s (...); 또는 model에 // @dto 선언", tableName),
		})
	}

	return errs
}
