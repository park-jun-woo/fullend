//ff:func feature=ssac-validate type=util
//ff:what @call input value에서 Go 타입을 결정한다
package validator

import (
	"strings"

	"github.com/jinzhu/inflection"
)

// resolveCallInputType는 @call input value에서 Go 타입을 결정한다.
func resolveCallInputType(val string, resultModels map[string]string, st *SymbolTable) string {
	// 리터럴
	if strings.HasPrefix(val, `"`) {
		return "string"
	}

	dotIdx := strings.IndexByte(val, '.')
	if dotIdx < 0 {
		return ""
	}
	source := val[:dotIdx]
	field := val[dotIdx+1:]

	// currentUser → 현재는 타입 추적 불가, 스킵
	if source == "currentUser" {
		return ""
	}
	// request → DDL에서 역추적
	if source == "request" {
		snakeName := toSnakeCase(field)
		for _, table := range st.DDLTables {
			if goType, ok := table.Columns[snakeName]; ok {
				return goType
			}
		}
		return ""
	}
	// 변수.Field → 해당 변수의 모델 테이블에서 Field 컬럼 타입
	modelName, ok := resultModels[source]
	if !ok {
		return ""
	}
	tableName := inflection.Plural(toSnakeCase(modelName))
	snakeName := toSnakeCase(field)
	if table, ok := st.DDLTables[tableName]; ok {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	// ID 패턴 fallback
	if strings.HasSuffix(field, "ID") {
		refModel := field[:len(field)-2]
		refTable := inflection.Plural(toSnakeCase(refModel))
		if table, ok := st.DDLTables[refTable]; ok {
			if goType, ok := table.Columns["id"]; ok {
				return goType
			}
		}
	}
	return ""
}
