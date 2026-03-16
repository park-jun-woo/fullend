//ff:func feature=ssac-gen type=util control=sequence
//ff:what 모델명에 해당하는 DDL 테이블에서 컬럼 타입을 조회
package generator

import (
	"github.com/jinzhu/inflection"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func lookupModelTableColumn(modelName, snakeName string, st *validator.SymbolTable) string {
	tableName := inflection.Plural(toSnakeCase(modelName))
	if table, ok := st.DDLTables[tableName]; ok {
		if goType, ok := table.Columns[snakeName]; ok {
			return goType
		}
	}
	return ""
}
