//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=string-convert
//ff:what PascalCase/camelCase를 snake_case로 변환
package generator

// toSnakeCase는 PascalCase/camelCase를 snake_case로 변환한다.
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		result = appendSnakeChar(result, s, i, c)
	}
	return string(result)
}
