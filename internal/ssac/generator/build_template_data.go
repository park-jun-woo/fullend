//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what 시퀀스를 분석하여 templateData를 구성 (모델 호출, 가드, 상태, 인증, 응답 등)
package generator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func buildTemplateData(seq parser.Sequence, errDeclared *bool, declaredVars map[string]bool, resultTypes map[string]string, st *validator.SymbolTable, funcName string, useTx bool, resolver *FieldTypeResolver) templateData {
	d := templateData{
		Message: seq.Message,
		Result:  seq.Result,
	}

	buildModelCall(&d, seq, useTx, st)
	buildDefaultMessage(&d, seq)
	buildArgsForData(&d, seq, st)
	buildHasTotal(&d, seq)
	buildGuard(&d, seq, resolver, resultTypes)
	buildStateAuth(&d, seq)
	buildErrStatusForGuard(&d, seq)
	buildInputFieldsForData(&d, seq)
	buildPublishData(&d, seq)
	buildResponseData(&d, seq)
	buildResultTracking(&d, seq, declaredVars)
	buildErrTracking(&d, seq, errDeclared)

	return d
}
