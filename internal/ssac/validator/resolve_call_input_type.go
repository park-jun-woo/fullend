//ff:func feature=ssac-validate type=util control=sequence topic=type-resolve
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
		return lookupAnyColumnType(st, toSnakeCase(field))
	}
	// 변수.Field → 해당 변수의 모델 테이블에서 Field 컬럼 타입
	modelName, ok := resultModels[source]
	if !ok {
		return ""
	}
	tableName := inflection.Plural(toSnakeCase(modelName))
	snakeName := toSnakeCase(field)
	if goType := lookupColumnType(st, tableName, snakeName); goType != "" {
		return goType
	}
	// ID 패턴 fallback
	if !strings.HasSuffix(field, "ID") {
		return ""
	}
	refModel := field[:len(field)-2]
	refTable := inflection.Plural(toSnakeCase(refModel))
	return lookupColumnType(st, refTable, "id")
}
