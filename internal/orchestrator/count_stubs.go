//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what FuncSpec 목록에서 TODO 스텁 수를 집계
package orchestrator

import "github.com/park-jun-woo/fullend/internal/funcspec"

func countStubs(specs []funcspec.FuncSpec) int {
	count := 0
	for _, s := range specs {
		if !s.HasBody {
			count++
		}
	}
	return count
}
