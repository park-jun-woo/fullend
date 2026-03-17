//ff:func feature=crosscheck type=util control=iteration dimension=2 topic=openapi-ddl
//ff:what DDL 테이블에서 컬럼의 VARCHAR 길이와 CHECK enum을 조회
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

func findDDLColumnConstraints(st *ssacvalidator.SymbolTable, col string) (varcharLen int, checkEnums []string, found bool) {
	for _, table := range st.DDLTables {
		if _, ok := table.Columns[col]; !ok {
			continue
		}
		found = true
		if table.VarcharLen != nil {
			if n, ok := table.VarcharLen[col]; ok {
				varcharLen = n
			}
		}
		if table.CheckEnums != nil {
			if vals, ok := table.CheckEnums[col]; ok {
				checkEnums = vals
			}
		}
		return
	}
	return
}
