//ff:func feature=crosscheck type=rule control=sequence
//ff:what CheckAuthzCustomPackageSkips: 커스텀 authz 패키지 사용 시 검증을 건너뛰는지 테스트
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckAuthzCustomPackageSkips(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "DoSomething",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"CustomField": "value",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "github.com/custom/authz")
	if len(errs) != 0 {
		t.Fatalf("expected no errors with custom package, got %d", len(errs))
	}
}
