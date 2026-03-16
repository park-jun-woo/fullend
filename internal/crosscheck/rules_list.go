//ff:func feature=crosscheck type=util control=sequence
//ff:what 등록된 교차 검증 규칙 목록을 반환
package crosscheck

// Rules returns the registered rule list (for status/reporting).
func Rules() []Rule {
	return rules
}
