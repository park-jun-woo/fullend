//ff:func feature=stml-validate type=util control=iteration dimension=1
//ff:what 문자열 슬라이스에 특정 값이 포함되어 있는지 확인
package validator

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
