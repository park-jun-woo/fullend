//ff:func feature=crosscheck type=rule control=sequence
//ff:what 단일 x-include 스펙 항목의 포맷과 DDL 참조 유효성 검증
package crosscheck

import (
	"fmt"
	"strings"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// validateXIncludeSpec validates a single x-include spec entry.
func validateXIncludeSpec(spec, ctx, primaryTable string, st *ssacvalidator.SymbolTable) []CrossError {
	colonIdx := strings.Index(spec, ":")
	if colonIdx <= 0 {
		return []CrossError{{
			Rule:       "x-include ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("x-include %q: invalid format, expected 'column:table.column'", spec),
			Suggestion: "예시: instructor_id:users.id",
		}}
	}
	localCol := spec[:colonIdx]
	targetRef := spec[colonIdx+1:]
	dotIdx := strings.Index(targetRef, ".")
	if dotIdx <= 0 {
		return []CrossError{{
			Rule:       "x-include ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("x-include %q: invalid format, expected 'column:table.column'", spec),
			Suggestion: "예시: instructor_id:users.id",
		}}
	}
	targetTable := targetRef[:dotIdx]

	if _, ok := st.DDLTables[targetTable]; !ok {
		return []CrossError{{
			Rule:       "x-include ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("x-include %q: target table %q not found in DDL", spec, targetTable),
			Suggestion: fmt.Sprintf("DDL에 추가: CREATE TABLE %s (...);", targetTable),
		}}
	}

	if primaryTable != "" && !hasFKColumn(primaryTable, localCol, targetTable, st) {
		return []CrossError{{
			Rule:       "x-include ↔ DDL FK",
			Context:    ctx,
			Message:    fmt.Sprintf("x-include %q: column %s.%s does not reference %s", spec, primaryTable, localCol, targetTable),
			Level:      "WARNING",
			Suggestion: fmt.Sprintf("DDL에 추가: ALTER TABLE %s ADD COLUMN %s BIGINT REFERENCES %s(id);", primaryTable, localCol, targetTable),
		}}
	}

	return nil
}
