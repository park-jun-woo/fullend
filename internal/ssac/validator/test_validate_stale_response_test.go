//ff:func feature=ssac-validate type=test control=sequence
//ff:what stale 데이터 response 사용 시 WARNING 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidateStaleResponse(t *testing.T) {
	funcs := []parser.ServiceFunc{{
		Name: "Cancel", FileName: "cancel.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Reservation.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Reservation", Var: "reservation"}},
			{Type: parser.SeqPut, Model: "Reservation.UpdateStatus", Inputs: map[string]string{"ID": "request.ID", "Status": `"cancelled"`}},
			{Type: parser.SeqResponse, Fields: map[string]string{"reservation": "reservation"}},
		},
	}}
	errs := Validate(funcs)
	assertHasWarning(t, errs, "갱신 없이 response에 사용됩니다")
}
