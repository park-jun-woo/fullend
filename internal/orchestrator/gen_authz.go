//ff:func feature=orchestrator type=command control=sequence
//ff:what genAuthz generates OPA authorizer package from parsed Rego policies (pkg 경로).

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/generate/gogin"
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
)

func genAuthz(artifactsDir string, policies []rego.Policy) reporter.StepResult {
	step := reporter.StepResult{Name: "authz-gen"}

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
	step.Summary = fmt.Sprintf("OPA authorizer generated (%d rules)", countRegoPolicyRules(policies))
	return step
}
