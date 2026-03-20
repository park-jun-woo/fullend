//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what RunRulesSkipAll: 모든 규칙을 건너뛸 때 에러가 0인지 테스트
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/policy"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestRunRules_SkipAll(t *testing.T) {
	input := &CrossValidateInput{
		ParsedSSOTs: &genapi.ParsedSSOTs{
			ServiceFuncs: []ssacparser.ServiceFunc{
				{Name: "DummyFunc"},
			},
			Policies: []*policy.Policy{{
				File:  "test.rego",
				Rules: []policy.AllowRule{{Actions: []string{"Do"}, Resource: "r", RoleValue: "admin", SourceLine: 1}},
			}},
		},
		Roles: []string{"admin"},
	}

	// Skip all rules by name.
	skip := make(map[string]bool)
	for _, r := range Rules() {
		skip[r.Name] = true
	}

	errs := RunRules(input, skip)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors when all rules skipped, got %d", len(errs))
	}
}
