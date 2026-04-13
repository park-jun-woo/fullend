//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what source.field 참조에서 DDL 테이블의 컬럼 타입을 조회
package ssac

import (
	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func resolveSourceFieldType(source, field string, st *rule.Ground) string {
	if source == "" || source == "request" || source == "currentUser" {
		return ""
	}
	refTable := inflection.Plural(toSnakeCase(source))
	refCol := toSnakeCase(field)
	if table, ok := st.Tables[refTable]; ok {
		if goType, ok := table.Columns[refCol]; ok {
			return goType
		}
	}
	return ""
}
