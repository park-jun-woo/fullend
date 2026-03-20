//ff:func feature=ssac-validate type=test control=sequence
//ff:what package-prefix @model 사용 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidatePackagePrefixModelError(t *testing.T) {
	st := &SymbolTable{Models: map[string]ModelSymbol{}, Operations: map[string]OperationSymbol{}, DDLTables: map[string]DDLTable{}}
	funcs := []parser.ServiceFunc{{
		Name: "GetSession", FileName: "get_session.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Package: "session", Model: "Session.Get", Inputs: map[string]string{"token": "request.Token"}, Result: &parser.Result{Type: "Session", Var: "session"}},
		},
	}}
	errs := ValidateWithSymbols(funcs, st)
	assertHasError(t, errs, "package-prefix @model은 지원하지 않습니다")
	assertHasError(t, errs, "@call Func을 사용하세요")
}
