//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 소스 테이블에 특정 FK 컬럼이 대상 테이블을 참조하는지 확인
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// hasFKColumn checks if srcTable has a FK column named colName that references dstTable.
func hasFKColumn(srcTable, colName, dstTable string, st *ssacvalidator.SymbolTable) bool {
	table, ok := st.DDLTables[srcTable]
	if !ok {
		return false
	}
	for _, fk := range table.ForeignKeys {
		if fk.Column == colName && fk.RefTable == dstTable {
			return true
		}
	}
	return false
}
