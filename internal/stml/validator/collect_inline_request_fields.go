//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 인라인 스키마의 properties에서 요청 필드를 수집
package validator

func collectInlineRequestFields(schema openAPISchema, fields map[string]string) {
	for fname, fprop := range schema.Properties {
		fields[fname] = fprop.Type
	}
}
