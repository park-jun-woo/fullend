//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectFieldKeys — map[string]string의 키를 슬라이스로 수집
package crosscheck

func collectFieldKeys(fields map[string]string) []string {
	var keys []string
	for k := range fields {
		keys = append(keys, k)
	}
	return keys
}
