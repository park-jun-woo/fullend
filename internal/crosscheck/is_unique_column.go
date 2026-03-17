//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what 컬럼이 PRIMARY KEY 또는 UNIQUE 제약인지 확인
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

// isUniqueColumn checks if a column is PRIMARY KEY or has a UNIQUE constraint.
func isUniqueColumn(col, tableName string, st *ssacvalidator.SymbolTable) bool {
	table, ok := st.DDLTables[tableName]
	if !ok {
		return false
	}
	for _, pk := range table.PrimaryKey {
		if pk == col {
			return true
		}
	}
	for _, idx := range table.Indexes {
		if idx.IsUnique && len(idx.Columns) == 1 && idx.Columns[0] == col {
			return true
		}
	}
	return false
}
