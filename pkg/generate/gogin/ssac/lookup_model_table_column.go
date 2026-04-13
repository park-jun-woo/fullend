//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what 모델명에 해당하는 DDL 테이블에서 컬럼 타입을 조회
package ssac

import (
	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func lookupModelTableColumn(modelName, snakeName string, st *rule.Ground) string {
	tableName := inflection.Plural(toSnakeCase(modelName))
	if table, ok := st.Tables[tableName]; ok {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return ""
}
