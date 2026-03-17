//ff:func feature=ssac-gen type=generator control=selection topic=args-inputs
//ff:what 시퀀스 타입에 따라 Args 또는 Inputs를 코드로 변환
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildArgsForData(d *templateData, seq parser.Sequence, st *validator.SymbolTable) {
	switch seq.Type {
	case parser.SeqGet, parser.SeqPost, parser.SeqPut, parser.SeqDelete:
		var paramOrder []string
		if st != nil {
			paramOrder = lookupParamOrder(seq.Model, st)
		}
		d.ArgsCode = buildArgsCodeFromInputs(seq.Inputs, paramOrder)
	default:
		d.ArgsCode = buildArgsCode(seq.Args)
	}
}
