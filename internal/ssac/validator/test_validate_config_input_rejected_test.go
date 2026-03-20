//ff:func feature=ssac-validate type=test control=sequence
//ff:what config.* 입력 거부 ERROR 검증
package validator

import ("testing"; "github.com/park-jun-woo/fullend/internal/ssac/parser")

func TestValidateConfigInputRejected(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "SendEmail", FileName: "send_email.go", Sequences: []parser.Sequence{
		{Type: parser.SeqCall, Model: "mail.Send", Inputs: map[string]string{"Host": "config.SMTPHost"}},
	}}}
	errs := Validate(funcs)
	assertHasError(t, errs, "config.* 입력은 지원하지 않습니다")
}
