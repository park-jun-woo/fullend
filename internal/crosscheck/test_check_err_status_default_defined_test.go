//ff:func feature=crosscheck type=test control=sequence topic=ssac-openapi
//ff:what checkErrStatus: 기본 404 응답 정의 시 에러 없음 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckErrStatus_DefaultDefined(t *testing.T) {
	// @empty with default 404, OpenAPI has 404 response → no error.
	doc := buildErrStatusDoc("GetGig", "404")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "gig.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:   "empty",
			Target: "gig",
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}
