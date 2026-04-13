//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what DDL 테이블 전체에서 파라미터명에 해당하는 Go 타입을 조회
package ssac

import "github.com/park-jun-woo/fullend/pkg/rule"

func lookupDDLType(paramName string, st *rule.Ground) string {
	snakeName := toSnakeCase(paramName)
	for _, table := range st.Tables {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return "string"
}
