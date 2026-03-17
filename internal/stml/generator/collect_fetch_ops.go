//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what FetchBlock의 OperationID를 재귀적으로 수집한다
package generator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func collectFetchOps(f parser.FetchBlock, ops []string) []string {
	ops = append(ops, f.OperationID)
	for _, child := range f.NestedFetches {
		ops = collectFetchOps(child, ops)
	}
	return ops
}
