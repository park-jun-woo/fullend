//ff:func feature=gen-gogin type=util
//ff:what extracts a string value from a map with a default fallback

package gogin

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
