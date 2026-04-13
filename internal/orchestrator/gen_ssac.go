//ff:func feature=orchestrator type=command control=sequence
//ff:what genSSaC generates service functions and model interfaces from SSaC specs (pkg 경로).

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	ssacgenerator "github.com/park-jun-woo/fullend/pkg/generate/gogin/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func genSSaC(profile *TargetProfile, specsDir, artifactsDir string, fs *fullend.Fullstack, g *rule.Ground) []reporter.StepResult {
	var steps []reporter.StepResult

	funcs := fs.ServiceFuncs
	if funcs == nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{"SSaC parse failed"},
		})
		return steps
	}

	if g == nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{"Ground not available"},
		})
		return steps
	}

	if gt, ok := profile.Backend.(*ssacgenerator.GoTarget); ok {
		gt.FuncSpecs = append(fs.FullendPkgSpecs, fs.ProjectFuncSpecs...)
	}

	serviceOutDir := filepath.Join(artifactsDir, "backend", "internal", "service")
	if err := os.MkdirAll(serviceOutDir, 0755); err != nil {
		steps = append(steps, reporter.StepResult{
			Name:   "ssac-gen",
			Status: reporter.Fail,
			Errors: []string{fmt.Sprintf("cannot create dir: %v", err)},
		})
		return steps
	}

	step := reporter.StepResult{Name: "ssac-gen"}
	if err := ssacgenerator.GenerateWith(profile.Backend, funcs, serviceOutDir, g); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("SSaC generate error: %v", err))
	} else {
		step.Status = reporter.Pass
		step.Summary = fmt.Sprintf("%d service files generated", len(funcs))
	}
	steps = append(steps, step)

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
	if err := profile.Backend.GenerateModelInterfaces(funcs, g, modelOutDir); err != nil {
		modelStep.Status = reporter.Fail
		modelStep.Errors = append(modelStep.Errors, fmt.Sprintf("SSaC model interface error: %v", err))
	} else {
		modelStep.Status = reporter.Pass
		modelStep.Summary = "model interfaces generated"
	}
	steps = append(steps, modelStep)

	return steps
}
