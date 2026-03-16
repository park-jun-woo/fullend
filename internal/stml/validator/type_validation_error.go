//ff:type feature=stml-validate type=model
//ff:what STML 검증 오류를 나타내는 구조체
package validator

// ValidationError represents a single validation failure.
type ValidationError struct {
	File    string // source HTML filename
	Attr    string // attribute context (e.g. `data-fetch="Login"`)
	Message string // human-readable error
}
