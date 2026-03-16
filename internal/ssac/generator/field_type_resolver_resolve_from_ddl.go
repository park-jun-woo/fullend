//ff:func feature=ssac-gen type=util control=sequence
//ff:what DDL 테이블에서 모델명+필드명으로 Go 타입을 조회
package generator

import "github.com/jinzhu/inflection"

func (r *FieldTypeResolver) resolveFromDDL(modelName, fieldName string) string {
	tableName := inflection.Plural(toSnakeCase(modelName))
	table, ok := r.st.DDLTables[tableName]
	if !ok {
		return ""
	}
	snakeField := toSnakeCase(fieldName)
	if goType, ok := table.Columns[snakeField]; ok {
		return goType
	}
	return ""
}
