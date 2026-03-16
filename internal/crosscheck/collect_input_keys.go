//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what map의 키를 bool set으로 변환
package crosscheck

func collectInputKeys(inputs map[string]string) map[string]bool {
	fields := make(map[string]bool, len(inputs))
	for k := range inputs {
		fields[k] = true
	}
	return fields
}
