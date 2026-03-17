//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what snake_case 컬럼이 DDL 테이블 중 하나에 존재하는지 확인
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

// columnExistsInAnyTable checks if a snake_case column exists in any DDL table.
func columnExistsInAnyTable(snake string, st *ssacvalidator.SymbolTable) bool {
	for _, table := range st.DDLTables {
		if _, ok := table.Columns[snake]; ok {
			return true
		}
	}
	return false
}
