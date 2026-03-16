//ff:func feature=orchestrator type=command control=sequence
//ff:what genSSaC generates service functions and model interfaces from SSaC specs.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/genapi"
	"github.com/geul-org/fullend/internal/reporter"
	ssacgenerator "github.com/geul-org/fullend/internal/ssac/generator"
)

func genSSaC(profile *TargetProfile, specsDir, artifactsDir string, parsed *genapi.ParsedSSOTs) []reporter.StepResult {
	var steps []reporter.StepResult

	funcs := parsed.ServiceFuncs
	if funcs == nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{"SSaC parse failed"},
		})
		return steps
	}

	if parsed.SymbolTable == nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{"SSaC symbol table not available"},
		})
		return steps
	}

	// Clone SymbolTable for gen path — injectFuncErrStatus mutates Models.
	genST := parsed.SymbolTable.Clone()

	// Inject @error annotations from func specs into the cloned symbol table.
	injectFuncErrStatusFromParsed(genST, parsed)

	// Generate service functions → backend/internal/service/
	serviceOutDir := filepath.Join(artifactsDir, "backend", "internal", "service")
	if err := os.MkdirAll(serviceOutDir, 0755); err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create dir: %v", err)},
		})
		return steps
	}

	if gt, ok := profile.Backend.(*ssacgenerator.GoTarget); ok {
		gt.FuncSpecs = append(parsed.FullendPkgSpecs, parsed.ProjectFuncSpecs...)
	}

	step := reporter.StepResult{Name: "ssac-gen"}
	if err := ssacgenerator.GenerateWith(profile.Backend, funcs, serviceOutDir, genST); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("SSaC generate error: %v", err))
	} else {
		step.Status = reporter.Pass
		step.Summary = fmt.Sprintf("%d service files generated", len(funcs))
	}
	steps = append(steps, step)

	// Generate model interfaces → backend/internal/model/
	// SSaC writes to outDir/model/, so pass backend/internal/ as outDir.
	modelOutDir := filepath.Join(artifactsDir, "backend", "internal")
	if err := os.MkdirAll(modelOutDir, 0755); err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-model",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create dir: %v", err)},
		})
		return steps
	}

	modelStep := reporter.StepResult{Name: "ssac-model"}
	if err := profile.Backend.GenerateModelInterfaces(funcs, genST, modelOutDir); err != nil {
		modelStep.Status = reporter.Fail
		modelStep.Errors = append(modelStep.Errors, fmt.Sprintf("SSaC model interface error: %v", err))
	} else {
		modelStep.Status = reporter.Pass
		modelStep.Summary = "model interfaces generated"
	}
	steps = append(steps, modelStep)

	return steps
}
