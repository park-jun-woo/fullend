//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 지정된 규칙을 건너뛰며 교차 검증 규칙을 실행
package crosscheck

// RunRules executes rules, skipping names in skipRules.
func RunRules(input *CrossValidateInput, skipRules map[string]bool) []CrossError {
	var errs []CrossError
	for _, r := range rules {
		if skipRules[r.Name] {
			continue
		}
		if r.Requires(input) {
			errs = append(errs, r.Check(input)...)
		}
	}
	return errs
}
