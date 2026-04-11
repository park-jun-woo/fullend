//ff:type feature=crosscheck type=model
//ff:what CrossError — 교차 검증 위반 결과
package crosscheck

// CrossError represents a single cross-validation violation.
type CrossError struct {
	Rule       string // rule ID: "X-1", "X-15", etc.
	Context    string // e.g. operationId, funcName, table name
	Message    string
	Level      string // "ERROR" or "WARNING" (empty = ERROR)
	Suggestion string
}
