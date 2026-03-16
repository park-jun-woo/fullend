//ff:func feature=gen-hurl type=util control=sequence
//ff:what map에서 기본값 포함 문자열을 추출한다
package hurl

// getStr extracts a string value from a map with a default.
func getStr(m map[string]interface{}, key, def string) string {
	v, ok := m[key]
	if !ok {
		return def
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return def
}
