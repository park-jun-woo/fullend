//ff:type feature=crosscheck type=model
//ff:what 교차 검증 규칙의 메타데이터와 실행 함수를 담는 구조체
package crosscheck

// Rule represents a single cross-validation rule with metadata.
type Rule struct {
	Name     string // e.g. "OpenAPI ↔ DDL", "SSaC → OpenAPI"
	Source   string // "OpenAPI", "SSaC", "Policy", "States", "Config", "DDL"
	Target   string // "DDL", "OpenAPI", ... ("" = standalone)
	Requires func(*CrossValidateInput) bool
	Check    func(*CrossValidateInput) []CrossError
}
