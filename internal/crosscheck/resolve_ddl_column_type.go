//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what SymbolTable에서 DDL 컬럼의 Go 타입 조회
package crosscheck

import (
	"strings"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// resolveDDLColumnType looks up a column's Go type from the SymbolTable.
func resolveDDLColumnType(st *ssacvalidator.SymbolTable, tableName, columnName string) string {
	if st == nil || st.DDLTables == nil {
		return ""
	}
	table, ok := st.DDLTables[tableName]
	if !ok {
		return ""
	}
	if goType, ok := table.Columns[columnName]; ok {
		return goType
	}
	snakeCol := toSnakeCase(columnName)
	if goType, ok := table.Columns[snakeCol]; ok {
		return goType
	}
	for colName, goType := range table.Columns {
		if strings.EqualFold(colName, columnName) {
			return goType
		}
	}
	return ""
}
