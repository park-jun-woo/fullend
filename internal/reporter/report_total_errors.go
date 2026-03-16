//ff:func feature=reporter type=util control=iteration dimension=1
//ff:what 모든 단계의 에러 수 합계를 반환한다
package reporter

// TotalErrors returns the sum of errors across all steps.
func (r *Report) TotalErrors() int {
	n := 0
	for _, s := range r.Steps {
		n += len(s.Errors)
	}
	return n
}
