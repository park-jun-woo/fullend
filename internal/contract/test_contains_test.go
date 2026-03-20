//ff:func feature=contract type=util control=sequence
//ff:what contains: 문자열 포함 여부를 확인하는 테스트 헬퍼
package contract

func contains(s, sub string) bool {
	return len(s) >= len(sub) && containsStr(s, sub)
}
