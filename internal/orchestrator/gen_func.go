//ff:func feature=orchestrator type=command control=iteration dimension=2
//ff:what genFunc copies custom func spec Go files to the artifacts directory.

package orchestrator

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/reporter"
)

func genFunc(funcDir, specsDir, artifactsDir, modulePath string) reporter.StepResult {
	step := reporter.StepResult{Name: "func-gen"}

	entries, err := os.ReadDir(funcDir)
	if err != nil {
		step.Status = reporter.Skip
		step.Summary = "no func/ directory"
		return step
	}

	// Scan SSaC files to find import paths for each func package.
	funcImportPaths, err := scanFuncImports(specsDir, modulePath)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("failed to scan SSaC imports: %v", err))
		return step
	}

	copied := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		n := copyFuncPackage(entry.Name(), funcDir, artifactsDir, modulePath, funcImportPaths, &step)
		if n < 0 {
			return step
		}
		copied += n
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d func files copied", copied)
	return step
}
