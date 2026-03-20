//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_ResponseIgnoredWarning: func spec에 Response 있지만 @call에 result 없으면 WARNING 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_ResponseIgnoredWarning(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package:        "auth",
		Name:           "doSomething",
		RequestFields:  []funcspec.Field{},
		ResponseFields: []funcspec.Field{{Name: "Value", Type: "string"}},
		HasBody:        true,
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:   "call",
			Model:  "auth.DoSomething",
			Inputs: nil,
			Result: nil, // no result
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "WARNING" && contains(e.Message, "반환값 무시") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected response ignored WARNING, got: %+v", errs)
	}
}
