//ff:func feature=gen-hurl type=util control=iteration
//ff:what 응답 스키마를 추출한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// getResponseSchema extracts the 2xx success response schema from an operation.
func getResponseSchema(op *openapi3.Operation) *openapi3.Schema {
	if op.Responses == nil {
		return nil
	}
	// Try explicit 2xx codes first, then fall back to 200.
	for code, respRef := range op.Responses.Map() {
		if len(code) == 3 && code[0] == '2' && respRef != nil && respRef.Value != nil && respRef.Value.Content != nil {
			ct := respRef.Value.Content.Get("application/json")
			if ct != nil && ct.Schema != nil {
				return resolveSchema(ct.Schema)
			}
		}
	}
	return nil
}
