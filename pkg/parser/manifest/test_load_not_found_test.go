//ff:func feature=manifest type=parser control=sequence
//ff:what fullend.yaml 미존재 시 에러 반환 검증
package manifest

import "testing"

func TestLoad_NotFound(t *testing.T) {
	_, err := Load("/nonexistent/dir")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
