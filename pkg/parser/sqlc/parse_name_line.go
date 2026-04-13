//ff:func feature=sqlc-parse type=util control=sequence
//ff:what parseNameLine — "-- name: FindByID :one" 라인에서 이름 + cardinality 추출
package sqlc

import "strings"

// parseNameLine은 "-- name: FindByID :one" 형식의 줄에서 메서드명과 cardinality 를 추출한다.
// modelName 접두가 붙어 있으면 제거한다 ("CourseFindByID" + "Course" → "FindByID").
func parseNameLine(line, modelName string) (name, cardinality string) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return "", ""
	}
	name = stripModelPrefix(parts[2], modelName)
	if len(parts) >= 4 {
		cardinality = strings.TrimPrefix(parts[3], ":")
	}
	return name, cardinality
}
