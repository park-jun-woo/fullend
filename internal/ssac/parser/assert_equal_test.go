//ff:func feature=ssac-parse type=util control=sequence
//ff:what assertEqual: 문자열 비교 실패 시 name·got·want 출력하는 테스트 헬퍼
package parser

import "testing"

func assertEqual(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", name, got, want)
	}
}
