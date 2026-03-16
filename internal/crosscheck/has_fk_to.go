//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what 소스 테이블에서 대상 테이블로의 FK 존재 확인
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// hasFKTo checks if srcTable has a FK pointing to dstTable.
func hasFKTo(srcTable, dstTable string, st *ssacvalidator.SymbolTable) bool {
	table, ok := st.DDLTables[srcTable]
	if !ok {
		return false
	}
	for _, fk := range table.ForeignKeys {
		if fk.RefTable == dstTable {
			return true
		}
	}
	return false
}
