//ff:func feature=orchestrator type=util
//ff:what traceFuncSpecs finds func specs referenced by @call sequences.

package orchestrator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func traceFuncSpecs(sf *ssacparser.ServiceFunc, specs []funcspec.FuncSpec, specsDir string) []ChainLink {
	callPkgFuncs := map[string]string{} // "pkg.Func" -> pkg
	for _, seq := range sf.Sequences {
		if seq.Type != "call" || seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) == 2 {
			callPkgFuncs[seq.Model] = parts[0]
		}
	}

	if len(callPkgFuncs) == 0 {
		return nil
	}

	var links []ChainLink
	for callRef, pkg := range callPkgFuncs {
		parts := strings.SplitN(callRef, ".", 2)
		funcName := ""
		if len(parts) == 2 {
			funcName = parts[1]
		}
		for _, spec := range specs {
			if spec.Package == pkg && strings.EqualFold(spec.Name, funcName) {
				relPath := "func/" + spec.Package + "/" + toSnakeCase(spec.Name) + ".go"
				// Try to find actual file.
				if _, err := os.Stat(filepath.Join(specsDir, relPath)); err != nil {
					// Try glob.
					matches, _ := filepath.Glob(filepath.Join(specsDir, "func", spec.Package, "*.go"))
					for _, m := range matches {
						if grepLine(m, "@func") > 0 && grepLine(m, funcName) > 0 {
							rel, _ := filepath.Rel(specsDir, m)
							relPath = rel
							break
						}
					}
				}
				line := grepLine(filepath.Join(specsDir, relPath), funcName)
				links = append(links, ChainLink{
					Kind:    "FuncSpec",
					File:    relPath,
					Line:    line,
					Summary: "@func " + callRef,
				})
				break
			}
		}
	}
	return links
}
