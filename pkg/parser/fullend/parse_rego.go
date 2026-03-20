//ff:func feature=orchestrator type=parser control=iteration dimension=1
//ff:what 디렉토리 내 모든 .rego 파일을 opa/ast로 파싱
package fullend

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/open-policy-agent/opa/ast"
)

// parseRegoDir parses all .rego files in the given directory using opa/ast.
func parseRegoDir(dir string) []*ast.Module {
	var modules []*ast.Module
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".rego") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		module, err := ast.ParseModule(e.Name(), string(data))
		if err != nil {
			continue
		}
		modules = append(modules, module)
	}
	return modules
}
