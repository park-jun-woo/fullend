//ff:func feature=orchestrator type=command
//ff:what GenWith runs code generation with the specified TargetProfile.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/reporter"
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

	// 2. sqlc generate (exec) — auto-generate sqlc.yaml if needed.
	if _, ok := has[KindDDL]; ok {
		report.Steps = append(report.Steps, genSqlc(specsDir, artifactsDir))
	}

	// 3. oapi-codegen (exec) → backend/internal/api/
	if _, ok := has[KindOpenAPI]; ok {
		report.Steps = append(report.Steps, genOpenAPI(specsDir, artifactsDir))
	}

	// 4. SSaC Generate → backend/internal/service/
	// 5. SSaC GenerateModelInterfaces → backend/internal/model/
	if _, ok := has[KindSSaC]; ok {
		report.Steps = append(report.Steps, genSSaC(profile, specsDir, artifactsDir, parsed)...)
	}

	// 6. STML Generate → frontend/src/pages/
	var stmlDeps map[string]string
	var stmlPages []string
	var stmlPageOps map[string]string
	if _, ok := has[KindSTML]; ok {
		var step reporter.StepResult
		step, stmlDeps, stmlPages, stmlPageOps = genSTML(profile, specsDir, artifactsDir, parsed.STMLPages)
		report.Steps = append(report.Steps, step)
	}

	// 7. Glue code generation (Server struct + main.go + frontend setup)
	report.Steps = append(report.Steps, genGlue(specsDir, artifactsDir, has, parsed, stmlDeps, stmlPages, stmlPageOps))

	// 8. Hurl smoke test generation (part of glue-gen, report separately)
	{
		testsDir := filepath.Join(artifactsDir, "tests")
		if _, err := os.Stat(filepath.Join(testsDir, "smoke.hurl")); err == nil {
			report.Steps = append(report.Steps, reporter.StepResult{
				Name:    "hurl-gen",
				Status:  reporter.Pass,
				Summary: "smoke.hurl generated",
			})
		}
	}

	// 9. State machine code generation.
	if _, ok := has[KindStates]; ok {
		report.Steps = append(report.Steps, genStateMachines(specsDir, artifactsDir, parsed))
	}

	// 10. OPA Authorizer code generation.
	if _, ok := has[KindPolicy]; ok {
		report.Steps = append(report.Steps, genAuthz(artifactsDir, parsed))
	}

	// 11. Scenario: user writes .hurl directly, no generation needed.

	// 12. Func copy (custom func specs → artifacts).
	if d, ok := has[KindFunc]; ok {
		modulePath := determineModulePath(specsDir, artifactsDir, parsed.Config)
		report.Steps = append(report.Steps, genFunc(d.Path, specsDir, artifactsDir, modulePath))
	}

	// Post-gen: restore preserved function bodies.
	if preserveSnap != nil {
		hasPreserves := len(preserveSnap.FilePreserves) > 0 || len(preserveSnap.FuncPreserves) > 0
		if hasPreserves {
			warnings := contract.RestorePreserved(preserveSnap)
			preserveCount := len(preserveSnap.FilePreserves)
			for _, funcs := range preserveSnap.FuncPreserves {
				preserveCount += len(funcs)
			}
			step := reporter.StepResult{
				Name:    "preserve",
				Status:  reporter.Pass,
				Summary: fmt.Sprintf("%d preserved", preserveCount),
			}
			for _, w := range warnings {
				step.Errors = append(step.Errors,
					fmt.Sprintf("[WARN] contract changed: %s:%s (old=%s new=%s)",
						w.File, w.Function, w.OldContract, w.NewContract))
			}
			if len(warnings) > 0 {
				step.Suggestions = append(step.Suggestions, ".new 파일에서 새 계약 코드를 확인하세요")
			}
			report.Steps = append(report.Steps, step)
		}
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
