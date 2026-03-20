//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what SuppressWarn으로 stale response WARNING 억제 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateStaleResponseSuppressed(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ID": "request.ID", "Status": `"cancelled"`}},
			{Type: parser.SeqResponse, Fields: map[string]string{"reservation": "reservation"}, SuppressWarn: true},
		},
	}}
	errs := Validate(funcs)
	for _, e := range errs {
		if e.IsWarning() && contains(e.Message, "갱신 없이") {
			t.Errorf("expected stale WARNING to be suppressed: %s", e.Message)
		}
	}
}
