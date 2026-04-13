//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what 전체 DDL 테이블을 순회하여 컬럼 타입을 조회
package ssac

import "github.com/park-jun-woo/fullend/pkg/rule"

func lookupAllTablesColumn(snakeName string, st *rule.Ground) string {
	for _, table := range st.Tables {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return "string"
}
