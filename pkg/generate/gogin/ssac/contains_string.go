//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=string-convert
//ff:what 문자열 슬라이스에 대상 문자열이 포함되어 있는지 확인
package ssac

func containsString(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
