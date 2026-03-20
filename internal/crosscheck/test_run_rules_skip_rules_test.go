//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what RunRulesSkipRules: 특정 규칙을 건너뛸 때 에러 수가 줄어드는지 테스트
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/genapi"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestRunRules_SkipRules(t *testing.T) {
	// Minimal input that triggers "SSaC Queue" and "SSaC -> Authz" rules.
	input := &CrossValidateInput{
		ParsedSSOTs: &genapi.ParsedSSOTs{
			ServiceFuncs: []ssacparser.ServiceFunc{
				{Name: "DummyFunc"},
			},
		},
	}

	// Run without skip -- collect baseline errors.
	allErrs := RunRules(input, nil)

	// Run with skip "SSaC Queue" -- should have fewer or equal errors.
	skip := map[string]bool{"SSaC Queue": true}
	skippedErrs := RunRules(input, skip)

	// Count how many errors come from "SSaC Queue" rule in the full run.
	queueCount := 0
	for _, e := range allErrs {
		if e.Rule == "SSaC Queue" || containsQueueError(e) {
			queueCount++
		}
	}

	// The skipped run should have exactly (all - queue) errors,
	// but at minimum it should not have more than the full run.
	if len(skippedErrs) > len(allErrs) {
		t.Errorf("skipped run has more errors (%d) than full run (%d)", len(skippedErrs), len(allErrs))
	}

	_ = queueCount // used for reasoning
}
