//ff:func feature=ssac-validate type=util control=iteration dimension=2 topic=type-resolve
//ff:what snake_case 컬럼명이 존재하는 DDL 테이블을 찾는다
package validator

import (
	"strings"

	"github.com/jinzhu/inflection"
)

// findColumnTable는 snake_case 컬럼명이 존재하는 DDL 테이블을 찾는다.
func findColumnTable(snakeCol, model string, st *SymbolTable) (string, bool) {
	if st == nil {
		return "", false
	}
	// 모델명에서 테이블명 유추: "Transaction.Create" → "transactions"
	if model != "" {
		parts := strings.SplitN(model, ".", 2)
		tableName := inflection.Plural(toSnakeCase(parts[0]))
		table, ok := st.DDLTables[tableName]
		if ok {
			if _, ok := table.Columns[snakeCol]; ok {
				return tableName, true
			}
		}
	}
	// 전체 DDL 테이블에서 검색
	for tableName, table := range st.DDLTables {
		if _, ok := table.Columns[snakeCol]; ok {
			return tableName, true
		}
	}
	return "", false
}
