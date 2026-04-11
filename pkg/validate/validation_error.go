//ff:type feature=rule type=model
//ff:what ValidationError — 단일 SSOT 검증 위반 결과
package validate

// ValidationError represents a single validation violation.
type ValidationError struct {
	Rule    string // rule ID: "S-1", "TM-1", "C-1", etc.
	File    string
	Func    string // SSaC func name or page name
	SeqIdx  int    // sequence index (-1 if N/A)
	Level   string // "ERROR" or "WARNING"
	Message string
}
