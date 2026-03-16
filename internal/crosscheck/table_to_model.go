//ff:func feature=crosscheck type=util control=sequence
//ff:what DDL 테이블 이름을 모델 이름으로 변환
package crosscheck

import "github.com/jinzhu/inflection"

// tableToModel converts a DDL table name to a model name.
// "courses" → "Course", "enrollments" → "Enrollment"
func tableToModel(table string) string {
	return snakeToPascal(inflection.Singular(table))
}
