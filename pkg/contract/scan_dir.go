//ff:func feature=contract type=walker control=sequence
//ff:what 아티팩트 디렉토리의 Go 파일에서 fullend 디렉티브를 스캔한다
package contract

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ScanDir scans artifacts directory for all Go files with //fullend: directives.
func ScanDir(artifactsDir string) ([]FuncStatus, error) {
	var results []FuncStatus

	backendDir := filepath.Join(artifactsDir, "backend")
	err := filepath.WalkDir(backendDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(artifactsDir, path)
		content := string(src)

		// Check file-level directive.
		if d := findFileLevelDirective(content); d != nil {
			results = append(results, FuncStatus{
				File:      relPath,
				Function:  "(file)",
				Directive: *d,
				Status:    d.Ownership,
			})
			return nil
		}

		// Check function-level directives.
		funcStatuses := scanFuncDirectives(content, relPath)
		results = append(results, funcStatuses...)
		return nil
	})

	return results, err
}
