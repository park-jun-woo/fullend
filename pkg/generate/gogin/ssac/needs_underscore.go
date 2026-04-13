//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what snake_case 변환 시 언더스코어 삽입이 필요한지 판단
package ssac

func needsUnderscore(s string, i int) bool {
	prev := s[i-1]
	if prev >= 'a' && prev <= 'z' {
		return true
	}
	if prev >= 'A' && prev <= 'Z' && i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z' {
		return true
	}
	return false
}
