//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what map에서 문자열 슬라이스를 추출한다
package hurl

// getStrSlice extracts a string slice from a map.
func getStrSlice(m map[string]interface{}, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
