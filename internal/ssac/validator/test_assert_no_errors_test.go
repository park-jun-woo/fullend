//ff:func feature=ssac-validate type=test-helper control=iteration dimension=1
//ff:what ERROR 레벨 검증 에러가 없는지 확인하는 테스트 헬퍼
package validator

import "testing"

func assertNoErrors(t *testing.T, errs []ValidationError) {
	t.Helper()
	var errors []ValidationError
	for _, e := range errs {
		if e.Level != "WARNING" {
			errors = append(errors, e)
		}
	}
	if len(errors) > 0 {
		t.Errorf("expected no errors, got %d: %v", len(errors), errors)
	}
}
