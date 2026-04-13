//ff:func feature=ssac-gen type=test-helper control=sequence
//ff:what 코드에 substring이 포함되지 않았는지 검증하는 테스트 헬퍼
package ssac

import (
	"strings"
	"testing"
)

func assertNotContains(t *testing.T, code, substr string) {
	t.Helper()
	if strings.Contains(code, substr) {
		t.Errorf("expected code NOT to contain %q\n--- code ---\n%s", substr, code)
	}
}
