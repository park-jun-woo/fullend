//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what Inputs value에서 DDL 테이블을 참조하여 Go 타입을 추론
package ssac

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

// resolveInputParamType는 Inputs value에서 Go 타입을 추론한다.
// value 형식: "request.Field", "source.Field", "\"literal\"", "currentUser.Field"
func resolveInputParamType(val string, modelName string, st *rule.Ground) string {
	if strings.HasPrefix(val, `"`) {
		return "string"
	}

	dotIdx := strings.IndexByte(val, '.')
	if dotIdx < 0 {
		return "string"
	}
	source := val[:dotIdx]
	field := val[dotIdx+1:]

	if source != "request" && source != "currentUser" {
		if goType := resolveSourceFieldType(source, field, st); goType != "" {
			return goType
		}
	}

	snakeName := toSnakeCase(field)

	tableName := inflection.Plural(toSnakeCase(modelName))
	if table, ok := st.Tables[tableName]; ok {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}

	if goType := lookupRefIDType(field, st); goType != "" {
		return goType
	}

	return lookupAllTablesColumn(snakeName, st)
}
