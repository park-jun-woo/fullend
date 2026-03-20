//ff:func feature=crosscheck type=test control=sequence topic=ssac-openapi
//ff:what checkErrStatus: 커스텀 402 응답 정의 시 에러 없음 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckErrStatus_CustomDefined(t *testing.T) {
	// @empty with custom 402, OpenAPI has 402 response → no error.
	doc := buildErrStatusDoc("ExecuteWorkflow", "402")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "ExecuteWorkflow",
		FileName: "workflow.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:      "empty",
			Target:    "org",
			ErrStatus: 402,
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}
