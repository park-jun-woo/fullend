//ff:func feature=sqlc-parse type=util control=sequence
//ff:what stripModelPrefix — 쿼리 이름에서 모델명 접두사 제거
package sqlc

import "strings"

// stripModelPrefix는 쿼리 이름에서 모델명 접두사를 제거한다.
// "CourseFindByID" + "Course" → "FindByID"
// "FindByID" + "Course" → "FindByID"
func stripModelPrefix(queryName, modelName string) string {
	if !strings.HasPrefix(queryName, modelName) {
		return queryName
	}
	stripped := queryName[len(modelName):]
	if len(stripped) > 0 && stripped[0] >= 'A' && stripped[0] <= 'Z' {
		return stripped
	}
	return queryName
}
