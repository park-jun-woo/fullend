//ff:type feature=sqlc-parse type=model
//ff:what sqlc SQL 쿼리 파싱 결과 타입
package sqlc

// Query는 sqlc 쿼리 파일(`.sql`) 의 `-- name: Xxx :one` 항목 하나를 나타낸다.
type Query struct {
	Model       string   // 파일명 단수화 + PascalCase (e.g. "users.sql" → "User")
	Name        string   // 쿼리 이름 (모델 prefix 제거 후, e.g. "FindByEmail")
	Cardinality string   // "one" / "many" / "exec"
	Params      []string // $1, $2 순서의 PascalCase 파라미터 이름
}
