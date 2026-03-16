//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 지정 컬럼을 포함하는 첫 번째 테이블명 반환
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// findTableWithColumn returns the first table name containing the given column.
func findTableWithColumn(col string, st *ssacvalidator.SymbolTable) string {
	for tableName, table := range st.DDLTables {
		if _, ok := table.Columns[col]; ok {
			return tableName
		}
	}
	return "???"
}
