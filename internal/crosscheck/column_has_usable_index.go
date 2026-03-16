//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what 컬럼에 사용 가능한 인덱스(선행 컬럼 또는 단일 컬럼)가 있는지 확인
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// columnHasUsableIndex checks if a column has a usable index.
func columnHasUsableIndex(col string, st *ssacvalidator.SymbolTable) bool {
	for _, table := range st.DDLTables {
		if _, ok := table.Columns[col]; !ok {
			continue
		}
		if tableHasLeadingIndex(col, table) {
			return true
		}
	}
	return false
}
