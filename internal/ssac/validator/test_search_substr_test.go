//ff:func feature=ssac-validate type=test-helper control=iteration dimension=1
//ff:what 문자열에서 부분 문자열을 순차 탐색하는 테스트 헬퍼
package validator

func searchSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
