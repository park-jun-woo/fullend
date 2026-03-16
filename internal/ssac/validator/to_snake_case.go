//ff:func feature=ssac-validate type=util control=iteration dimension=1
//ff:what PascalCase/camelCase를 snake_case로 변환한다
package validator

// toSnakeCase는 PascalCase/camelCase를 snake_case로 변환한다.
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c < 'A' || c > 'Z' {
			result = append(result, byte(c))
			continue
		}
		if i > 0 && needsUnderscore(s, i) {
			result = append(result, '_')
		}
		result = append(result, byte(c)+32)
	}
	return string(result)
}
