//ff:type feature=genapi type=model
//ff:what STML 생성기 출력 결과를 보관하는 타입
package genapi

// STMLGenOutput holds STML generator output (not parse results).
// Populated by orchestrator after stml.Generate(), consumed by react gen.
type STMLGenOutput struct {
	Deps    map[string]string // npm dependencies
	Pages   []string          // page names
	PageOps map[string]string // page file → primary operationID
}
