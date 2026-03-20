//ff:func feature=crosscheck type=util control=sequence
//ff:what contains: 문자열 포함 여부를 확인하는 테스트 헬퍼
package crosscheck

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}
