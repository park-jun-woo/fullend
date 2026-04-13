//ff:func feature=rule type=loader control=iteration dimension=2
//ff:what populateOps — OpenAPI operation 메타를 g.Ops 로 복사
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateOps(g *rule.Ground, fs *fullend.Fullstack) {
	if g.Ops == nil {
		g.Ops = make(map[string]rule.OperationInfo)
	}
	if fs.OpenAPIDoc == nil {
		return
	}
	for path, pathItem := range fs.OpenAPIDoc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			if op.OperationID == "" {
				continue
			}
			g.Ops[op.OperationID] = rule.OperationInfo{
				ID:             op.OperationID,
				Method:         method,
				Path:           path,
				PathParams:     extractPathParams(op.Parameters),
				HasRequestBody: op.RequestBody != nil,
				Pagination:     extractPagination(op),
				Sort:           extractSort(op),
				Filter:         extractFilter(op),
			}
		}
	}
}
