//ff:func feature=stml-validate type=test-helper control=iteration dimension=1
//ff:what 특정 문자열을 포함하는 에러가 있는지 확인하는 테스트 헬퍼
package validator

import ("strings"; "testing")

func assertHasError(t *testing.T, errs []ValidationError, substr string) {
	t.Helper()
	for _, e := range errs { if strings.Contains(e.Error(), substr) { return } }
	t.Errorf("expected error containing %q, got: %v", substr, errs)
}
