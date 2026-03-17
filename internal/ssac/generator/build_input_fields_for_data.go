//ff:func feature=ssac-gen type=generator control=sequence topic=args-inputs
//ff:what state/auth/call 시퀀스의 Inputs를 Go struct 필드로 변환
package generator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

func buildInputFieldsForData(d *templateData, seq parser.Sequence) {
	if seq.Type != parser.SeqState && seq.Type != parser.SeqAuth && seq.Type != parser.SeqCall {
		return
	}
	if len(seq.Inputs) == 0 {
		return
	}
	inputs := seq.Inputs
	if seq.Type == parser.SeqAuth {
		inputs = filterAuthInputs(seq.Inputs)
	}
	d.InputFields = buildInputFieldsFromMap(inputs)
}
