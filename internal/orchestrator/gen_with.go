//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what GenWith runs code generation with the specified TargetProfile.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/contract"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

// GenWith runs code generation with the specified TargetProfile.
// When reset is false (default), preserved function bodies are restored after generation.
func GenWith(profile *TargetProfile, specsDir, artifactsDir string, skipKinds map[SSOTKind]bool, reset ...bool) (*reporter.Report, bool) {
	isReset := len(reset) > 0 && reset[0]
	detected, err := DetectSSOTs(specsDir)
	if err != nil {
		report := &reporter.Report{}
		report.Steps = append(report.Steps, reporter.StepResult{
			Name:   "detect",
			Status: reporter.Fail,
			Errors: []string{err.Error()},
		})
		return report, false
	}

	skip := skipKinds
	if skip == nil {
		skip = make(map[SSOTKind]bool)
	}

	// Parse all SSOTs once.
	parsed := ParseAll(specsDir, detected, skip)

	// 1. Validate first (reuse parsed data).
	report := ValidateWith(specsDir, detected, parsed, skip)
	if report.HasFailure() {
		return report, false
	}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	// Add separator between validate and gen steps.
	report.Steps = append(report.Steps, reporter.StepResult{
		Name:    "---",
		Status:  reporter.Pass,
		Summary: "codegen",
	})

	// Ensure artifacts directory exists.
	if err := os.MkdirAll(artifactsDir, 0755); err != nil {
		report.Steps = append(report.Steps, reporter.StepResult{
			Name:   "gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create artifacts dir: %v", err)},
		})
		return report, false
	}

	// Pre-gen: scan preserved functions (unless reset).
	var preserveSnap *contract.PreserveSnapshot
	backendDir := filepath.Join(artifactsDir, "backend")
	if !isReset {
		preserveSnap = contract.ScanPreserveSnapshot(backendDir)
	}

	// 2-12. Run all codegen steps.
	runCodegenSteps(report, profile, specsDir, artifactsDir, has, parsed)

	// Post-gen: restore preserved function bodies.
	if step := restorePreserved(preserveSnap); step != nil {
		report.Steps = append(report.Steps, *step)
	}

	genOk := true
	for _, s := range report.Steps {
		if s.Status == reporter.Fail {
			genOk = false
			break
		}
	}

	return report, genOk
}
