//ff:func feature=crosscheck type=test control=sequence topic=ssac-openapi
//ff:what checkErrStatus: 기본 404 응답 미정의 시 에러 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckErrStatus_DefaultMissing(t *testing.T) {
	// @empty with default 404, OpenAPI has no 404 response → error.
	doc := buildErrStatusDoc("GetGig", "200")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "gig.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:   "empty",
			Target: "gig",
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !contains(errs[0].Message, "404") {
		t.Errorf("expected 404 in message, got: %s", errs[0].Message)
	}
}
