//ff:func feature=orchestrator type=rule control=sequence
//ff:what OPA Rego 정책 검증 — pkg/validate/rego 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func validatePolicy(policies []*policy.Policy) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindPolicy)}
	if policies == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "Policy parse failed")
		return step
	}
	if len(policies) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no policy files found"
		return step
	}

	totalRules := countPolicyRules(policies)
	totalOwnerships := countPolicyOwnerships(policies)
	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d files, %d rules, %d ownership mappings", len(policies), totalRules, totalOwnerships)
	return step
}
