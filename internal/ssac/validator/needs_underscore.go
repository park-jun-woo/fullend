//ff:func feature=ssac-validate type=util control=sequence
//ff:what snake_case 변환 시 언더스코어 삽입 여부를 판단한다
package validator

func needsUnderscore(s string, i int) bool {
	prev := s[i-1]
	if prev >= 'a' && prev <= 'z' {
		return true
	}
	return prev >= 'A' && prev <= 'Z' && i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z'
}
