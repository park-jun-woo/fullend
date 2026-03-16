//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what DDL 테이블이 SSaC에서 참조되는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckDDLCoverage validates that DDL tables and columns are referenced by SSaC/OpenAPI.
func CheckDDLCoverage(
	st *ssacvalidator.SymbolTable,
	funcs []ssacparser.ServiceFunc,
	archived *ArchivedInfo,
) []CrossError {
	var errs []CrossError

	if st == nil || len(st.DDLTables) == 0 {
		return errs
	}

	if archived == nil {
		archived = &ArchivedInfo{
			Tables:  make(map[string]bool),
			Columns: make(map[string]map[string]bool),
		}
	}

	// Build set of tables referenced by SSaC (@model and @result).
	referencedTables := buildReferencedTables(funcs)

	// Rule 1: DDL table → SSaC reference.
	for tableName := range st.DDLTables {
		if archived.Tables[tableName] {
			continue
		}
		if !referencedTables[tableName] {
			errs = append(errs, CrossError{
				Rule:       "DDL → SSaC",
				Context:    tableName,
				Message:    fmt.Sprintf("DDL 테이블 %q가 SSaC에서 참조되지 않습니다", tableName),
				Level:      "ERROR",
				Suggestion: "더 이상 사용하지 않는 테이블이면 DDL에 -- @archived를 추가하세요",
			})
		}
	}

	return errs
}
