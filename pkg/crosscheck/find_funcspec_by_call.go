//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what findFuncSpecByCall — "pkg.Func" 형태 @call Model 문자열에서 funcspec 조회

package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
)

// findFuncSpecByCall locates the FuncSpec for a @call's Model string like "billing.CheckCredits".
// Returns nil if not found. Delegates per-list matching to matchFuncSpec.
func findFuncSpecByCall(model string, specs ...[]funcspec.FuncSpec) *funcspec.FuncSpec {
	idx := strings.LastIndex(model, ".")
	if idx < 0 {
		return nil
	}
	pkg := model[:idx]
	name := model[idx+1:]
	for _, list := range specs {
		if fs := matchFuncSpec(list, pkg, name); fs != nil {
			return fs
		}
	}
	return nil
}
