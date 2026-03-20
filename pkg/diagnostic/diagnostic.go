//ff:type feature=orchestrator type=model
//ff:what 파싱/검증/교차검증 단계의 진단 메시지
package diagnostic

// Diagnostic represents a single diagnostic message from any phase.
type Diagnostic struct {
	File    string // source file path
	Line    int    // line number (0 if unknown)
	Phase   Phase  // parse, validate, crosscheck
	Level   Level  // error, warning
	Message string // human-readable message
	Ref     *Loc   // counterpart location (crosscheck only, nil otherwise)
}
