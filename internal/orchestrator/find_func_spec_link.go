//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what callRef에 매칭되는 FuncSpec을 찾아 ChainLink를 반환한다

package orchestrator

import (
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/funcspec"
)

// findFuncSpecLink searches specs for a matching func spec and returns a ChainLink.
func findFuncSpecLink(callRef, pkg, funcName string, specs []funcspec.FuncSpec, specsDir string) (ChainLink, bool) {
	for _, spec := range specs {
		if spec.Package != pkg || !strings.EqualFold(spec.Name, funcName) {
			continue
		}
		relPath := resolveFuncSpecPath(spec, funcName, specsDir)
		line := grepLine(filepath.Join(specsDir, relPath), funcName)
		return ChainLink{
			Kind:    "FuncSpec",
			File:    relPath,
			Line:    line,
			Summary: "@func " + callRef,
		}, true
	}
	return ChainLink{}, false
}
