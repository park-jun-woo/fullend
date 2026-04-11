//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populatePathOps — path의 각 operation에서 operationId, method, response 추출
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populatePathOps(g *rule.Ground, opIDs, methods rule.StringSet, ops map[string]*openapi3.Operation) {
	for method, op := range ops {
		methods[method] = true
		if op.OperationID != "" {
			opIDs[op.OperationID] = true
			populateResponseSchema(g, op.OperationID, op)
		}
	}
}
