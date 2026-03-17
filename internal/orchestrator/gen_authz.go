//ff:func feature=orchestrator type=command control=sequence
//ff:what genAuthz generates OPA authorizer package from parsed Rego policies.

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/gen/gogin"
	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func genAuthz(artifactsDir string, parsed *genapi.ParsedSSOTs) reporter.StepResult {
	step := reporter.StepResult{Name: "authz-gen"}

	policies := parsed.Policies
	if policies == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "Policy parse failed")
		return step
	}
	if len(policies) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no policy files"
		return step
	}

	if err := gogin.GenerateAuthzPackage(policies, artifactsDir); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("authz-gen error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("OPA authorizer generated (%d rules)", countPolicyRules(policies))
	return step
}
