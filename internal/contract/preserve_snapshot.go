//ff:type feature=contract type=model
//ff:what 코드 생성 전 보존 대상 전체를 캡처하는 구조체
package contract

// PreserveSnapshot captures all preserved functions/files before code generation.
type PreserveSnapshot struct {
	FilePreserves map[string]string                    // path → saved whole file content
	FuncPreserves map[string]map[string]*PreservedFunc // path → funcName → preserved body
}
