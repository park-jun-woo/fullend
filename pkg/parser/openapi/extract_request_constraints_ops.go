//ff:func feature=manifest type=parser control=iteration dimension=1
//ff:what extractRequestConstraintsOps — 단일 PathItem의 operation별 requestBody 제약조건 추출
package openapi

import "github.com/getkin/kin-openapi/openapi3"

func extractRequestConstraintsOps(result map[string]map[string]FieldConstraint, item *openapi3.PathItem) {
	for _, op := range item.Operations() {
		if op.OperationID == "" || op.RequestBody == nil {
			continue
		}
		fields := extractBodyConstraints(op.RequestBody, op.OperationID)
		if len(fields) > 0 {
			result[op.OperationID] = fields
		}
	}
}
