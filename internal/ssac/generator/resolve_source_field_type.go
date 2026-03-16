//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what source.field 참조에서 DDL 테이블의 컬럼 타입을 조회
package generator

import (
	"github.com/jinzhu/inflection"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func resolveSourceFieldType(source, field string, st *validator.SymbolTable) string {
	if source == "" || source == "request" || source == "currentUser" {
		return ""
	}
	refTable := inflection.Plural(toSnakeCase(source))
	refCol := toSnakeCase(field)
	if table, ok := st.DDLTables[refTable]; ok {
		if goType, ok := table.Columns[refCol]; ok {
			return goType
		}
	}
	return ""
}
