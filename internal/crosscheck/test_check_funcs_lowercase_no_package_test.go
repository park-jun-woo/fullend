//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_LowercaseNoPackage: 패키지 없이 소문자 함수명이면 ERROR 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_LowercaseNoPackage(t *testing.T) {
	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "someFunc",
			Inputs: map[string]string{
				"ID": "request.ID",
			},
		}},
	}}

	errs := CheckFuncs(sfs, nil, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "lowercase") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ERROR for lowercase func name without package, got: %+v", errs)
	}
}
