//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=string-convert
//ff:what 문자열 슬라이스에서 빈 문자열을 제거
package ssac

func filterNonEmpty(parts []string) []string {
	var result []string
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
