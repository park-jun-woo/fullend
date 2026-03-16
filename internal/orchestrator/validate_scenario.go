//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what Hurl 시나리오 검증 — .feature 파일 차단 + scenario/invariant .hurl 파일 수집
package orchestrator

import (
	"fmt"
	"path/filepath"

	"github.com/geul-org/fullend/internal/reporter"
)

func validateScenarioHurl(testsDir string, specsRoot string) (reporter.StepResult, []string) {
	step := reporter.StepResult{Name: string(KindScenario)}

	// Check for deprecated .feature files in both scenario/ (old) and tests/ (current).
	var featureFiles []string
	for _, dir := range []string{filepath.Join(specsRoot, "scenario"), testsDir} {
		if matches, _ := filepath.Glob(filepath.Join(dir, "*.feature")); len(matches) > 0 {
			featureFiles = append(featureFiles, matches...)
		}
	}
	if len(featureFiles) > 0 {
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
