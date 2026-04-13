//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what 교차 검증 실행 — pkg/crosscheck toulmin 기반 교차 정합성 검증
package orchestrator

import (
	"fmt"

	pkgcross "github.com/park-jun-woo/fullend/pkg/crosscheck"
	"github.com/park-jun-woo/fullend/pkg/fullend"

	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func runCrossValidate(root string, parsed *genapi.ParsedSSOTs) reporter.StepResult {
	step := reporter.StepResult{Name: "Cross"}

	if parsed.OpenAPIDoc == nil || parsed.SymbolTable == nil || parsed.ServiceFuncs == nil {
		step.Status = reporter.Skip
		step.Summary = "skipped (incomplete SSOT)"
		return step
	}

	detected, err := fullend.DetectSSOTs(root)
	if err != nil {
		step.Status = reporter.Fail
		step.Summary = "detect failed: " + err.Error()
		return step
	}

	fs := fullend.ParseAll(root, detected, nil)
	cerrs := pkgcross.Run(fs)

	hasError := false
	for _, ce := range cerrs {
		prefix := ce.Rule
		if ce.Level == "WARNING" {
			prefix = "[WARN] " + prefix
		} else {
			hasError = true
		}
		step.Errors = append(step.Errors, fmt.Sprintf("%s: %s — %s", prefix, ce.Context, ce.Message))
		step.Suggestions = append(step.Suggestions, ce.Suggestion)
	}

	if hasError {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}

	errCount, warnCount := countCrossErrors(cerrs)
	step.Summary = formatCrossSummary(errCount, warnCount)
	return step
}
