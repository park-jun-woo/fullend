//ff:func feature=ssac-validate type=test-helper control=iteration dimension=1
//ff:what 특정 문자열을 포함하는 WARNING이 있는지 확인하는 테스트 헬퍼
package validator

import "testing"

func assertHasWarning(t *testing.T, errs []ValidationError, substr string) {
	t.Helper()
	for _, e := range errs {
		if e.IsWarning() && contains(e.Message, substr) {
			return
		}
	}
	t.Errorf("expected warning containing %q, got %v", substr, errs)
}
