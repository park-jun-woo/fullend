//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what genSTML generates frontend pages from STML specs.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/reporter"
	stmlgenerator "github.com/geul-org/fullend/internal/stml/generator"
	stmlparser "github.com/geul-org/fullend/internal/stml/parser"
)

func genSTML(profile *TargetProfile, specsDir, artifactsDir string, pages []stmlparser.PageSpec) (reporter.StepResult, map[string]string, []string, map[string]string) {
	step := reporter.StepResult{Name: "stml-gen"}

	if pages == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "STML parse failed")
		return step, nil, nil, nil
	}

	// Output to frontend/src/pages/
	outDir := filepath.Join(artifactsDir, "frontend", "src", "pages")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir: %v", err))
		return step, nil, nil, nil
	}

	result, err := stmlgenerator.GenerateWith(profile.Frontend, pages, specsDir, outDir, stmlgenerator.GenerateOptions{
		APIImportPath: "../api",
		UseClient:     false,
	})
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("STML generate error: %v", err))
		return step, nil, nil, nil
	}

	// Collect generated page names and primary operationIDs for glue-gen.
	var pageNames []string
	pageOps := make(map[string]string)
	for _, p := range pages {
		pageNames = append(pageNames, p.Name)
		if len(p.Fetches) > 0 {
			pageOps[p.Name] = p.Fetches[0].OperationID
		} else if len(p.Actions) > 0 {
			pageOps[p.Name] = p.Actions[0].OperationID
		}
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d pages generated", result.Pages)
	return step, result.Dependencies, pageNames, pageOps
}
