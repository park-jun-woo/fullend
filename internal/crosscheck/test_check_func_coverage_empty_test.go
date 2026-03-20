//ff:func feature=crosscheck type=rule control=sequence topic=func-coverage
//ff:what TestCheckFuncCoverage_Empty: 빈 입력에서 에러 없음 확인
package crosscheck

import "testing"

func TestCheckFuncCoverage_Empty(t *testing.T) {
	errs := CheckFuncCoverage(nil, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errs))
	}
}
