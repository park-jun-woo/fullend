//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populatePathOpsParams — path의 각 operation에서 param/sort/filter 등록
package ground

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populatePathOpsParams(g *rule.Ground, ops map[string]*openapi3.Operation) {
	for _, op := range ops {
		if op.OperationID != "" {
			populateOpParams(g, op)
		}
	}
}
