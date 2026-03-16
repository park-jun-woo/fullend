//ff:func feature=ssac-gen type=util control=sequence
//ff:what 대문자면 언더스코어+소문자로, 아니면 그대로 바이트 슬라이스에 추가
package generator

func appendSnakeChar(result []byte, s string, i int, c rune) []byte {
	if c < 'A' || c > 'Z' {
		return append(result, byte(c))
	}
	if i > 0 && needsUnderscore(s, i) {
		result = append(result, '_')
	}
	return append(result, byte(c)+32)
}
