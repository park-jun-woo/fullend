//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_ForbiddenImport: func spec에 database/sql 금지 import가 있으면 ERROR 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_ForbiddenImport(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "bad",
		Name:    "doQuery",
		RequestFields: []funcspec.Field{
			{Name: "Key", Type: "string"},
		},
		HasBody: true,
		Imports: []string{"database/sql", "fmt"},
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "bad.DoQuery",
			Inputs: map[string]string{
				"Key": "request.Key",
			},
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	found := false
	for _, e := range errs {
		if e.Level == "ERROR" && contains(e.Message, "database/sql") && contains(e.Message, "I/O 패키지") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected forbidden import ERROR for database/sql, got: %+v", errs)
	}

	// fmt should NOT be flagged.
	for _, e := range errs {
		if contains(e.Message, `"fmt"`) {
			t.Errorf("fmt should not be forbidden: %+v", e)
		}
	}
}
