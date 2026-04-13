//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=string-convert
//ff:what 문자열의 선행 대문자 연속 개수를 반환
package ssac

func countLeadingUpper(s string) int {
	upper := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			upper++
		} else {
			break
		}
	}
	return upper
}
