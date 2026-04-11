//ff:func feature=orchestrator type=loader control=iteration dimension=1
//ff:what populateSTMLOps — operation 목록에서 operationId 수집
package orchestrator

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateSTMLOps(opIDs rule.StringSet, ops map[string]*openapi3.Operation) {
	for _, op := range ops {
		if op.OperationID != "" {
			opIDs[op.OperationID] = true
		}
	}
}
