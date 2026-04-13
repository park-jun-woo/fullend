//ff:func feature=ssac-gen type=generator control=selection topic=args-inputs
//ff:what 시퀀스 타입에 따라 Args 또는 Inputs를 코드로 변환
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func buildArgsForData(d *templateData, seq ssacparser.Sequence, st *rule.Ground) {
	switch seq.Type {
	case ssacparser.SeqGet, ssacparser.SeqPost, ssacparser.SeqPut, ssacparser.SeqDelete:
		var paramOrder []string
		if st != nil {
			paramOrder = lookupParamOrder(seq.Model, st)
		}
		d.ArgsCode = buildArgsCodeFromInputs(seq.Inputs, paramOrder)
	default:
		d.ArgsCode = buildArgsCode(seq.Args)
	}
}
