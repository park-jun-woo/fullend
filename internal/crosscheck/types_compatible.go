//ff:func feature=crosscheck type=util control=sequence topic=func-check
//ff:what 두 Go 타입 문자열의 호환성 확인
package crosscheck

// typesCompatible checks if two Go type strings are compatible.
func typesCompatible(a, b string) bool {
	return a == b
}
