//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what DDL 테이블 전체에서 파라미터명에 해당하는 Go 타입을 조회
package generator

import "github.com/geul-org/fullend/internal/ssac/validator"

func lookupDDLType(paramName string, st *validator.SymbolTable) string {
	snakeName := toSnakeCase(paramName)
	for _, table := range st.DDLTables {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return "string"
}
