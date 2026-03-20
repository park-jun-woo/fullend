//ff:func feature=contract type=util control=iteration dimension=1
//ff:what containsStr: 문자열 내 부분 문자열 검색을 수행하는 테스트 헬퍼
package contract

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
