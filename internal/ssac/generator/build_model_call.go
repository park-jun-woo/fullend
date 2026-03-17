//ff:func feature=ssac-gen type=generator control=selection topic=template-data
//ff:what 시퀀스의 모델 참조를 분석하여 templateData의 ModelCall/PkgName 설정
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildModelCall(d *templateData, seq parser.Sequence, useTx bool, st *validator.SymbolTable) {
	if seq.Model == "" {
		return
	}
	parts := strings.SplitN(seq.Model, ".", 2)
	switch seq.Type {
	case parser.SeqCall:
		buildCallModelData(d, parts, seq, st)
	default:
		buildCRUDModelData(d, parts, useTx)
	}
}
