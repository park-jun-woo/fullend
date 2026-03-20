//ff:func feature=crosscheck type=test control=sequence topic=ssac-openapi
//ff:what checkErrStatus: 커스텀 402 응답 미정의 시 에러 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckErrStatus_CustomMissing(t *testing.T) {
	// @empty with custom 402, OpenAPI has no 402 response → error.
	doc := buildErrStatusDoc("ExecuteWorkflow", "404")

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
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !contains(errs[0].Message, "402") {
		t.Errorf("expected 402 in message, got: %s", errs[0].Message)
	}
}
