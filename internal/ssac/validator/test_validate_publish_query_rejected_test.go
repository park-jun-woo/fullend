//ff:func feature=ssac-validate type=test control=sequence
//ff:what @publish에서 query 사용 시 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidatePublishQueryRejected(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "CreateOrder", FileName: "create_order.go", Sequences: []parser.Sequence{
		{Type: parser.SeqPublish, Topic: "order.created", Inputs: map[string]string{"Filter": "query"}},
	}}}
	errs := Validate(funcs)
	assertHasError(t, errs, "query는 HTTP 전용입니다")
}
