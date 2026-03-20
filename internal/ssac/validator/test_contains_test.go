//ff:func feature=ssac-validate type=test-helper control=sequence
//ff:what 문자열 포함 여부를 검사하는 테스트 헬퍼
package validator

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstr(s, substr)
}
