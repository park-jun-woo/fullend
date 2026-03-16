//ff:func feature=crosscheck type=util control=sequence
//ff:what 타입 이름에서 슬라이스·포인터 접두사를 제거
package crosscheck

// normalizeTypeName strips slice prefix and pointer prefix from a type name.
// e.g. "[]Reservation" → "Reservation", "*User" → "User"
func normalizeTypeName(t string) string {
	if len(t) > 2 && t[:2] == "[]" {
		t = t[2:]
	}
	if len(t) > 1 && t[0] == '*' {
		t = t[1:]
	}
	return t
}
