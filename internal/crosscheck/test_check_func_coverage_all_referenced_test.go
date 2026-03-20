//ff:func feature=crosscheck type=rule control=sequence topic=func-coverage
//ff:what TestCheckFuncCoverage_AllReferenced: 모든 func spec이 SSaC에서 참조되면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncCoverage_AllReferenced(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "CreateOrder",
			Sequences: []ssacparser.Sequence{
				{Type: "call", Model: "billing.HoldEscrow"},
			},
		},
	}
	specs := []funcspec.FuncSpec{
		{Package: "billing", Name: "HoldEscrow"},
	}

	errs := CheckFuncCoverage(funcs, specs)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}
