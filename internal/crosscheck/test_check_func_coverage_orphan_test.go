//ff:func feature=crosscheck type=rule control=sequence topic=func-coverage
//ff:what TestCheckFuncCoverage_Orphan: SSaC에서 참조되지 않는 func spec이 WARNING 발생 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncCoverage_Orphan(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name:      "CreateOrder",
			Sequences: []ssacparser.Sequence{},
		},
	}
	specs := []funcspec.FuncSpec{
		{Package: "billing", Name: "HoldEscrow"},
	}

	errs := CheckFuncCoverage(funcs, specs)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if errs[0].Level != "WARNING" {
		t.Errorf("expected WARNING, got %q", errs[0].Level)
	}
	if errs[0].Context != "billing.HoldEscrow" {
		t.Errorf("expected context billing.HoldEscrow, got %q", errs[0].Context)
	}
}
