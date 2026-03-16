//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 오퍼레이션의 200 응답 필드를 수집
package validator

// collectResponseFields extracts response fields from an operation's 200 response.
func collectResponseFields(op openAPIOperation, schemas map[string]openAPISchema, fields map[string]FieldSymbol) {
	resp, ok := op.Responses["200"]
	if !ok {
		return
	}
	for _, ct := range resp.Content {
		ref := ct.Schema.Ref
		if ref != "" {
			resolveResponseFields(schemas, ref, fields)
		} else {
			collectInlineResponseFields(ct.Schema, schemas, fields)
		}
	}
}
