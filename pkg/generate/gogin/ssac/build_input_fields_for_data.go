//ff:func feature=ssac-gen type=generator control=sequence topic=args-inputs
//ff:what state/auth/call 시퀀스의 Inputs를 Go struct 필드로 변환
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildInputFieldsForData(d *templateData, seq ssacparser.Sequence) {
	if seq.Type != ssacparser.SeqState && seq.Type != ssacparser.SeqAuth && seq.Type != ssacparser.SeqCall {
		return
	}
	if len(seq.Inputs) == 0 {
		return
	}
	inputs := seq.Inputs
	if seq.Type == ssacparser.SeqAuth {
		inputs = filterAuthInputs(seq.Inputs)
	}
	d.InputFields = buildInputFieldsFromMap(inputs)
}
