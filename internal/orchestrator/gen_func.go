//ff:func feature=orchestrator type=command control=iteration dimension=2
//ff:what genFunc copies custom func spec Go files to the artifacts directory.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/reporter"
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
		pkg := entry.Name()
		srcDir := filepath.Join(funcDir, pkg)

		// Determine destination from SSaC import path.
		importPath, ok := funcImportPaths[pkg]
		if !ok {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("func/%s: SSaC에서 import하는 곳이 없습니다", pkg))
			return step
		}

		// Extract relative path within module.
		relPath := strings.TrimPrefix(importPath, modulePath+"/")
		if relPath == importPath {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("func/%s: import 경로 %q가 모듈 %q에 속하지 않습니다", pkg, importPath, modulePath))
			return step
		}

		// Validate: must be under internal/ or pkg/.
		if !strings.HasPrefix(relPath, "internal/") && !strings.HasPrefix(relPath, "pkg/") {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("func/%s: import 경로 %q는 internal/ 또는 pkg/ 하위여야 합니다 (현재: %s)", pkg, importPath, relPath))
			return step
		}

		dstDir := filepath.Join(artifactsDir, "backend", relPath)

		if err := os.MkdirAll(dstDir, 0755); err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir %s: %v", dstDir, err))
			return step
		}

		files, _ := filepath.Glob(filepath.Join(srcDir, "*.go"))
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				step.Status = reporter.Fail
				step.Errors = append(step.Errors, fmt.Sprintf("read %s: %v", f, err))
				return step
			}
			dst := filepath.Join(dstDir, filepath.Base(f))
			if err := os.WriteFile(dst, data, 0644); err != nil {
				step.Status = reporter.Fail
				step.Errors = append(step.Errors, fmt.Sprintf("write %s: %v", dst, err))
				return step
			}
			copied++
		}
	}

	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("%d func files copied", copied)
	return step
}
