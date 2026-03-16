//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what 전체 DDL 테이블을 순회하여 컬럼 타입을 조회
package generator

import "github.com/geul-org/fullend/internal/ssac/validator"

func lookupAllTablesColumn(snakeName string, st *validator.SymbolTable) string {
	for _, table := range st.DDLTables {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return "string"
}
