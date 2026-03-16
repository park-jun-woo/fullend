//ff:type feature=crosscheck type=model
//ff:what 교차 검증 오류 정보를 담는 구조체
package crosscheck

// CrossError represents a cross-validation error between two SSOT layers.
type CrossError struct {
	Rule       string // e.g. "x-sort ↔ DDL", "SSaC @result ↔ DDL"
	Context    string // e.g. operationId or funcName
	Message    string
	Level      string // "ERROR" or "WARNING" (empty = ERROR)
	Suggestion string // fix suggestion (empty if none)
}
