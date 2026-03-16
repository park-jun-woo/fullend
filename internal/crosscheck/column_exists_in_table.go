//ff:func feature=crosscheck type=util control=sequence
//ff:what 특정 DDL 테이블에 컬럼이 존재하는지 확인
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// columnExistsInTable checks if a column exists in a specific DDL table.
func columnExistsInTable(st *ssacvalidator.SymbolTable, tableName, colName string) bool {
	if tbl, ok := st.DDLTables[tableName]; ok {
		if _, colOk := tbl.Columns[colName]; colOk {
			return true
		}
	}
	return false
}
