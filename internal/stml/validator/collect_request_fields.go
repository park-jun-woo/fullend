//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 오퍼레이션의 요청 본문 필드를 수집
package validator

// collectRequestFields extracts request body fields from an operation.
func collectRequestFields(op openAPIOperation, schemas map[string]openAPISchema, fields map[string]string) {
	if op.RequestBody.Content == nil {
		return
	}
	for _, ct := range op.RequestBody.Content {
		ref := ct.Schema.Ref
		if ref != "" {
			resolveSchemaFields(schemas, ref, fields)
		} else {
			collectInlineRequestFields(ct.Schema, fields)
		}
	}
}
