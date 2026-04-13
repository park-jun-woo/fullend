//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what {Model}ID 패턴의 필드에서 참조 테이블의 id 컬럼 타입을 조회
package ssac

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func lookupRefIDType(field string, st *rule.Ground) string {
	if !strings.HasSuffix(field, "ID") {
		return ""
	}
	refModel := field[:len(field)-2]
	refTable := inflection.Plural(toSnakeCase(refModel))
	if table, ok := st.Tables[refTable]; ok {
		if goType, ok := table.Columns["id"]; ok {
			return goType
		}
	}
	return ""
}
