//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 두 세그먼트 배열의 일치 여부 확인
package crosscheck

// segmentsMatch checks if two segment arrays match.
func segmentsMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
