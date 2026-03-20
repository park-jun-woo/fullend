//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_AllowedImportsOnly: 허용된 import만 사용하면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckFuncs_AllowedImportsOnly(t *testing.T) {
	specs := []funcspec.FuncSpec{{
		Package: "calc",
		Name:    "calculate",
		RequestFields: []funcspec.Field{
			{Name: "Value", Type: "int64"},
		},
		HasBody: true,
		Imports: []string{"math", "strings", "fmt", "time", "encoding/json"},
	}}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "calc.Calculate",
			Inputs: map[string]string{
				"Value": "request.Value",
			},
		}},
	}}

	errs := CheckFuncs(sfs, specs, nil, nil, nil)
	for _, e := range errs {
		if contains(e.Message, "I/O 패키지") {
			t.Errorf("unexpected forbidden import error: %+v", e)
		}
	}
}
