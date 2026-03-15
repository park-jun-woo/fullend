//ff:func feature=gen-hurl type=util
//ff:what 문자열 슬라이스를 정렬한다
package hurl

// sortStringSlice returns a sorted copy of the string slice.
func sortStringSlice(ss []string) []string {
	result := make([]string, len(ss))
	copy(result, ss)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
