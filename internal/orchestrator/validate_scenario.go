//ff:func feature=orchestrator type=rule
//ff:what Hurl 시나리오 검증 — .feature 파일 차단 + scenario/invariant .hurl 파일 수집
package orchestrator

import (
	"fmt"
	"path/filepath"

	"github.com/geul-org/fullend/internal/reporter"
)

func validateScenarioHurl(testsDir string, specsRoot string) (reporter.StepResult, []string) {
	step := reporter.StepResult{Name: string(KindScenario)}

	// Check for deprecated .feature files anywhere under specs root.
	scenarioDir := filepath.Join(specsRoot, "scenario")
	if featureFiles, _ := filepath.Glob(filepath.Join(scenarioDir, "*.feature")); len(featureFiles) > 0 {
		step.Status = reporter.Fail
		for _, f := range featureFiles {
			rel, _ := filepath.Rel(specsRoot, f)
			step.Errors = append(step.Errors, fmt.Sprintf("%s: .feature is no longer supported. Delete this file.\n       Write scenario tests directly in Hurl format: tests/scenario-*.hurl\n       See: https://hurl.dev/docs/manual.html", rel))
		}
		return step, nil
	}

	// Collect scenario and invariant .hurl files.
	scenarioHurls, _ := filepath.Glob(filepath.Join(testsDir, "scenario-*.hurl"))
	invariantHurls, _ := filepath.Glob(filepath.Join(testsDir, "invariant-*.hurl"))
	allHurls := append(scenarioHurls, invariantHurls...)

	if len(allHurls) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no scenario .hurl files found"
		return step, nil
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d scenario hurl files", len(allHurls))
	return step, allHurls
}
