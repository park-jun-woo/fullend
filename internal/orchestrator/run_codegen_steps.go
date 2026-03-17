//ff:func feature=orchestrator type=command control=sequence
//ff:what runs all code generation steps and appends results to report

package orchestrator

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

// runCodegenSteps runs all code generation steps and appends results to report.
func runCodegenSteps(report *reporter.Report, profile *TargetProfile, specsDir, artifactsDir string, has map[SSOTKind]DetectedSSOT, parsed *genapi.ParsedSSOTs) {
	if _, ok := has[KindDDL]; ok {
		report.Steps = append(report.Steps, genSqlc(specsDir, artifactsDir))
	}
	if _, ok := has[KindOpenAPI]; ok {
		report.Steps = append(report.Steps, genOpenAPI(specsDir, artifactsDir))
	}
	if _, ok := has[KindSSaC]; ok {
		report.Steps = append(report.Steps, genSSaC(profile, specsDir, artifactsDir, parsed)...)
	}

	var stmlDeps map[string]string
	var stmlPages []string
	var stmlPageOps map[string]string
	if _, ok := has[KindSTML]; ok {
		var step reporter.StepResult
		step, stmlDeps, stmlPages, stmlPageOps = genSTML(profile, specsDir, artifactsDir, parsed.STMLPages)
		report.Steps = append(report.Steps, step)
	}

	report.Steps = append(report.Steps, genGlue(specsDir, artifactsDir, has, parsed, stmlDeps, stmlPages, stmlPageOps))

	testsDir := filepath.Join(artifactsDir, "tests")
	if _, err := os.Stat(filepath.Join(testsDir, "smoke.hurl")); err == nil {
		report.Steps = append(report.Steps, reporter.StepResult{
			Name: "hurl-gen", Status: reporter.Pass, Summary: "smoke.hurl generated",
		})
	}

	if _, ok := has[KindStates]; ok {
		report.Steps = append(report.Steps, genStateMachines(specsDir, artifactsDir, parsed))
	}
	if _, ok := has[KindPolicy]; ok {
		report.Steps = append(report.Steps, genAuthz(artifactsDir, parsed))
	}
	if d, ok := has[KindFunc]; ok {
		modulePath := determineModulePath(specsDir, artifactsDir, parsed.Config)
		report.Steps = append(report.Steps, genFunc(d.Path, specsDir, artifactsDir, modulePath))
	}
}
