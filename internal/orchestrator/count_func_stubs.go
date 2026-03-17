//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what counts func specs without a body (TODO stubs)

package orchestrator

import "github.com/park-jun-woo/fullend/internal/funcspec"

// countFuncStubs counts func specs without a body (TODO stubs).
func countFuncStubs(specs []funcspec.FuncSpec) int {
	stubs := 0
	for _, s := range specs {
		if !s.HasBody {
			stubs++
		}
	}
	return stubs
}
