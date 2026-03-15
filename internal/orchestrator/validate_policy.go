//ff:func feature=orchestrator type=rule
//ff:what OPA Rego 정책 검증 — 파일 수, 규칙 수, 소유권 매핑 수 집계
package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/policy"
	"github.com/geul-org/fullend/internal/reporter"
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

	totalRules := 0
	totalOwnerships := 0
	for _, p := range policies {
		totalRules += len(p.Rules)
		totalOwnerships += len(p.Ownerships)
	}
	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d files, %d rules, %d ownership mappings", len(policies), totalRules, totalOwnerships)
	return step
}
