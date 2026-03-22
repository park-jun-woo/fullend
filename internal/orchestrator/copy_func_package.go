//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what 단일 func 패키지의 Go 파일들을 artifacts 디렉토리로 복사한다

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/internal/reporter"
)

// copyFuncPackage copies Go files from a single func package to the artifacts directory.
// Returns the number of files copied, or sets step error and returns -1 on failure.
func copyFuncPackage(pkg, funcDir, artifactsDir, modulePath string, funcImportPaths map[string]string, step *reporter.StepResult) int {
	srcDir := filepath.Join(funcDir, pkg)

	importPath, ok := funcImportPaths[pkg]
	if !ok {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("func/%s: SSaC에서 import하는 곳이 없습니다", pkg))
		return -1
	}

	relPath := strings.TrimPrefix(importPath, modulePath+"/")
	if relPath == importPath {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("func/%s: import 경로 %q가 모듈 %q에 속하지 않습니다", pkg, importPath, modulePath))
		return -1
	}

	if !strings.HasPrefix(relPath, "internal/") && !strings.HasPrefix(relPath, "pkg/") {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("func/%s: import 경로 %q는 internal/ 또는 pkg/ 하위여야 합니다 (현재: %s)", pkg, importPath, relPath))
		return -1
	}

	dstDir := filepath.Join(artifactsDir, "backend", relPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir %s: %v", dstDir, err))
		return -1
	}

	copied := 0
	files, _ := filepath.Glob(filepath.Join(srcDir, "*.go"))
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("read %s: %v", f, err))
			return -1
		}
		dst := filepath.Join(dstDir, filepath.Base(f))
		if err := os.WriteFile(dst, data, 0644); err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("write %s: %v", dst, err))
			return -1
		}
		copied++
	}
	return copied
}
