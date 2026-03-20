//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_ForbiddenImportNetHTTP: func spec에 net/http 금지 import가 있으면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_ForbiddenImportNetHTTP(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "bad",
		Name:    "callAPI",
		RequestFields: []funcspec.Field{
			{Name: "URL", Type: "string"},
		},
		HasBody: true,
		Imports: []string{"net/http"},
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "bad.CallAPI",
			Inputs: map[string]string{
				"URL": "request.URL",
			},
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "net/http") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected forbidden import ERROR for net/http, got: %+v", errs)
	}
}
