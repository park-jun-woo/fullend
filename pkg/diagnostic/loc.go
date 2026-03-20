//ff:type feature=orchestrator type=model
//ff:what 소스 위치 참조 (교차검증에서 상대 위치 지정)
package diagnostic

// Loc is a source location reference.
type Loc struct {
	File string
	Line int
}
