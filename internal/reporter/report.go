//ff:type feature=reporter type=model
//ff:what 전체 검증 결과를 담는 구조체
package reporter

// Report holds all step results from a validation run.
type Report struct {
	Steps []StepResult
}
