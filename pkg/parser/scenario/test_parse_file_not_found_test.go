//ff:func feature=scenario type=parser control=sequence
//ff:what 존재하지 않는 파일에 대해 nil 반환 검증
package scenario

import "testing"

func TestParseFile_NotFound(t *testing.T) {
	entries := ParseFile("/nonexistent/file.hurl")
	if entries != nil {
		t.Fatalf("expected nil, got %v", entries)
	}
}
