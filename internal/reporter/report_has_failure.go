//ff:func feature=reporter type=util control=iteration dimension=1
//ff:what 실패한 단계가 있는지 확인한다
package reporter

// HasFailure returns true if any step failed.
func (r *Report) HasFailure() bool {
	for _, s := range r.Steps {
		if s.Status == Fail {
			return true
		}
	}
	return false
}
