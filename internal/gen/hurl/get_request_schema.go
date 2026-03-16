//ff:func feature=gen-hurl type=util control=sequence
//ff:what 요청 본문 스키마를 추출한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// getRequestSchema extracts the request body schema from an operation.
func getRequestSchema(op *openapi3.Operation) *openapi3.Schema {
	if op.RequestBody == nil || op.RequestBody.Value == nil {
		return nil
	}
	ct := op.RequestBody.Value.Content.Get("application/json")
	if ct == nil || ct.Schema == nil {
		return nil
	}
	return resolveSchema(ct.Schema)
}
