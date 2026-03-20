//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what searchString: 문자열 내 부분 문자열 탐색
package crosscheck

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
