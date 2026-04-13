//ff:func feature=orchestrator type=command control=sequence
//ff:what runs all code generation steps and appends results to report (pkg 경로 통합)

package orchestrator

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

// runCodegenSteps runs all code generation steps and appends results to report.
func runCodegenSteps(report *reporter.Report, profile *TargetProfile, specsDir, artifactsDir string, has map[SSOTKind]DetectedSSOT) {
	if _, ok := has[KindDDL]; ok {
		report.Steps = append(report.Steps, genSqlc(specsDir, artifactsDir))
	}
	if _, ok := has[KindOpenAPI]; ok {
		report.Steps = append(report.Steps, genOpenAPI(specsDir, artifactsDir))
	}

	fs, g := buildPkgContext(specsDir)

	if _, ok := has[KindSSaC]; ok {
		report.Steps = append(report.Steps, genSSaC(profile, specsDir, artifactsDir, fs, g)...)
	}

	var stmlDeps map[string]string
	var stmlPages []string
	var stmlPageOps map[string]string
	if _, ok := has[KindSTML]; ok {
		var step reporter.StepResult
		step, stmlDeps, stmlPages, stmlPageOps = genSTML(profile, specsDir, artifactsDir, fs.STMLPages)
		report.Steps = append(report.Steps, step)
	}

	report.Steps = append(report.Steps, genGlue(specsDir, artifactsDir, fs, g, stmlDeps, stmlPages, stmlPageOps))

	testsDir := filepath.Join(artifactsDir, "tests")
	if _, err := os.Stat(filepath.Join(testsDir, "smoke.hurl")); err == nil {
		report.Steps = append(report.Steps, reporter.StepResult{
			Name: "hurl-gen", Status: reporter.Pass, Summary: "smoke.hurl generated",
		})
	}

	if _, ok := has[KindStates]; ok {
		report.Steps = append(report.Steps, genStateMachines(specsDir, artifactsDir, fs))
	}
	if d, ok := has[KindPolicy]; ok {
		policies, _ := policy.ParseDir(d.Path)
		report.Steps = append(report.Steps, genAuthz(artifactsDir, policies))
	}
	if d, ok := has[KindFunc]; ok {
		modulePath := determinePkgModulePath(specsDir, artifactsDir, fs.Manifest)
		report.Steps = append(report.Steps, genFunc(d.Path, specsDir, artifactsDir, modulePath))
	}
}
