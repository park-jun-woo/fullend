package crosscheck

import (
	"testing"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
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

func TestCheckFuncCoverage_Empty(t *testing.T) {
	errs := CheckFuncCoverage(nil, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errs))
	}
}
