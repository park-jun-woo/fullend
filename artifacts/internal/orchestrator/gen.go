package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	oapicodegen "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	oapiutil "github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	sqlccli "github.com/sqlc-dev/sqlc/pkg/cli"

	"github.com/geul-org/fullend/artifacts/internal/reporter"
	ssacgenerator "github.com/geul-org/ssac/generator"
	ssacparser "github.com/geul-org/ssac/parser"
	ssacvalidator "github.com/geul-org/ssac/validator"
	stmlgenerator "github.com/geul-org/stml/generator"
	stmlparser "github.com/geul-org/stml/parser"
)

// Gen runs validate first, then generates code from all detected SSOTs.
// Returns the validate report (with gen steps appended) and whether gen succeeded.
func Gen(specsDir, artifactsDir string) (*reporter.Report, bool) {
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

	// 1. Validate first.
	report := Validate(specsDir, detected)
	if report.HasFailure() {
		return report, false
	}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	// Terraform: optional — warn and skip if not installed.
	terraformAvailable := true
	if _, ok := has[KindTerraform]; ok {
		if _, err := exec.LookPath("terraform"); err != nil {
			terraformAvailable = false
		}
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

	// 2. sqlc generate (Go import)
	if _, ok := has[KindDDL]; ok {
		report.Steps = append(report.Steps, genSqlc(specsDir))
	}

	// 3. oapi-codegen (Go import)
	if _, ok := has[KindOpenAPI]; ok {
		report.Steps = append(report.Steps, genOpenAPI(specsDir, artifactsDir))
	}

	// 4. SSaC Generate (service functions)
	// 5. SSaC GenerateModelInterfaces
	if d, ok := has[KindSSaC]; ok {
		report.Steps = append(report.Steps, genSSaC(specsDir, d.Path, artifactsDir)...)
	}

	// 6. STML Generate (React TSX)
	if d, ok := has[KindSTML]; ok {
		report.Steps = append(report.Steps, genSTML(specsDir, d.Path, artifactsDir))
	}

	// 7. terraform fmt (외부 도구, 선택)
	if _, ok := has[KindTerraform]; ok {
		if terraformAvailable {
			report.Steps = append(report.Steps, genTerraform(specsDir))
		} else {
			report.Steps = append(report.Steps, reporter.StepResult{
				Name:    "terraform",
				Status:  reporter.Skip,
				Summary: "terraform 미설치, 스킵",
				Errors:  []string{"[WARN] terraform이 설치되어 있지 않습니다 — HCL 포맷팅을 건너뜁니다"},
			})
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

func genSqlc(specsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "sqlc"}
	configPath := filepath.Join(specsDir, "sqlc.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		step.Status = reporter.Skip
		step.Summary = "sqlc.yaml not found, skipped"
		return step
	}
	code := sqlccli.Run([]string{"generate", "-f", configPath})
	if code != 0 {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("sqlc generate failed (exit %d)", code))
		return step
	}
	step.Status = reporter.Pass
	step.Summary = "DB models generated"
	return step
}

func genOpenAPI(specsDir, artifactsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "oapi-gen"}
	apiPath := filepath.Join(specsDir, "api", "openapi.yaml")

	spec, err := oapiutil.LoadSwagger(apiPath)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("OpenAPI load error: %v", err))
		return step
	}

	outDir := filepath.Join(artifactsDir, "backend", "api")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir: %v", err))
		return step
	}

	// Generate types.
	typesCfg := oapicodegen.Configuration{
		PackageName: "api",
		Generate:    oapicodegen.GenerateOptions{Models: true},
	}
	typesCfg = typesCfg.UpdateDefaults()
	typesCode, err := oapicodegen.Generate(spec, typesCfg)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("oapi-codegen types error: %v", err))
		return step
	}
	if err := os.WriteFile(filepath.Join(outDir, "types.gen.go"), []byte(typesCode), 0644); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("write types.gen.go error: %v", err))
		return step
	}

	// Generate server (net/http std).
	serverCfg := oapicodegen.Configuration{
		PackageName: "api",
		Generate:    oapicodegen.GenerateOptions{StdHTTPServer: true},
	}
	serverCfg = serverCfg.UpdateDefaults()
	serverCode, err := oapicodegen.Generate(spec, serverCfg)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("oapi-codegen server error: %v", err))
		return step
	}
	if err := os.WriteFile(filepath.Join(outDir, "server.gen.go"), []byte(serverCode), 0644); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("write server.gen.go error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = "types + server generated"
	return step
}

func genSSaC(specsDir, serviceDir, artifactsDir string) []reporter.StepResult {
	var steps []reporter.StepResult

	funcs, err := ssacparser.ParseDir(serviceDir)
	if err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("SSaC parse error: %v", err)},
		})
		return steps
	}

	st, err := ssacvalidator.LoadSymbolTable(specsDir)
	if err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("SSaC symbol table error: %v", err)},
		})
		return steps
	}

	// Generate service functions.
	serviceOutDir := filepath.Join(artifactsDir, "backend", "service")
	if err := os.MkdirAll(serviceOutDir, 0755); err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create dir: %v", err)},
		})
		return steps
	}

	step := reporter.StepResult{Name: "ssac-gen"}
	if err := ssacgenerator.Generate(funcs, serviceOutDir, st); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("SSaC generate error: %v", err))
	} else {
		step.Status = reporter.Pass
		step.Summary = fmt.Sprintf("%d service files generated", len(funcs))
	}
	steps = append(steps, step)

	// Generate model interfaces.
	modelOutDir := filepath.Join(artifactsDir, "backend", "model")
	if err := os.MkdirAll(modelOutDir, 0755); err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-model",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create dir: %v", err)},
		})
		return steps
	}

	modelStep := reporter.StepResult{Name: "ssac-model"}
	if err := ssacgenerator.GenerateModelInterfaces(funcs, st, modelOutDir); err != nil {
		modelStep.Status = reporter.Fail
		modelStep.Errors = append(modelStep.Errors, fmt.Sprintf("SSaC model interface error: %v", err))
	} else {
		modelStep.Status = reporter.Pass
		modelStep.Summary = "model interfaces generated"
	}
	steps = append(steps, modelStep)

	return steps
}

func genSTML(specsDir, frontendDir, artifactsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "stml-gen"}

	pages, err := stmlparser.ParseDir(frontendDir)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("STML parse error: %v", err))
		return step
	}

	outDir := filepath.Join(artifactsDir, "frontend")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir: %v", err))
		return step
	}

	if err := stmlgenerator.Generate(pages, specsDir, outDir); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("STML generate error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d pages generated", len(pages))
	return step
}

func genTerraform(specsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "terraform"}
	tfDir := filepath.Join(specsDir, "terraform")
	res := RunExec("terraform", "fmt", tfDir)
	if res.Err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, res.Err.Error())
		if res.Stderr != "" {
			step.Errors = append(step.Errors, res.Stderr)
		}
		return step
	}
	step.Status = reporter.Pass
	step.Summary = "HCL formatted"
	return step
}
