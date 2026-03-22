//ff:func feature=contract type=test control=sequence
//ff:what 단일 Parse 테스트 케이스를 검증하는 헬퍼

package contract

import "testing"

// assertParseCase runs a single Parse test case.
func assertParseCase(t *testing.T, name, input string, want *Directive, wantErr bool) {
	t.Helper()
	got, err := Parse(input)
	if wantErr {
		if err == nil {
			t.Errorf("Parse(%q) expected error, got %+v", input, got)
		}
		return
	}
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", input, err)
	}
	if got.Ownership != want.Ownership || got.SSOT != want.SSOT || got.Contract != want.Contract {
		t.Errorf("Parse(%q) = %+v, want %+v", input, got, want)
	}
}
