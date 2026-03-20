//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what TestCheckFuncs_PackageModelSkipDDL: 패키지 접두사 모델(@get)은 CheckFuncs에서 @call 대상이 아님을 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckFuncs_PackageModelSkipDDL(t *testing.T) {
	// 패키지 접두사 모델 (@model session.Session.Get)은
	// DDL 테이블 체크 대상이 아님을 간접 확인.
	// CheckFuncs는 직접 DDL을 안 보지만, 패키지 모델이 @call이 아닌
	// @model인 경우 CheckSSaCDDL/CheckDDLCoverage에서 스킵됨.
	// 여기서는 패키지 모델이 specMap에 없어도 @call이 아니면 에러 안 남을 확인.

	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64"}},
		},
	}

	sfs := []ssacparser.ServiceFunc{{
		Name: "GetSession",
		Sequences: []ssacparser.Sequence{{
			Type:    "get",
			Package: "session",
			Model:   "Session.Get",
			Inputs:  map[string]string{"token": "request.Token"},
			Result:  &ssacparser.Result{Var: "session", Type: "Session"},
		}},
	}}

	// CheckFuncs only processes seq.Type == "call", so package model @get is ignored.
	errs := CheckFuncs(sfs, nil, nil, st, nil)
	for _, e := range errs {
		if e.Level == "ERROR" {
			t.Errorf("unexpected ERROR for package model: %+v", e)
		}
	}
}
