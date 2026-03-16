//ff:func feature=crosscheck type=util control=sequence
//ff:what 모든 교차 검증 규칙을 실행하고 결과를 반환
package crosscheck

// Run executes all cross-validation rules and returns collected errors.
func Run(input *CrossValidateInput) []CrossError {
	return RunRules(input, nil)
}
