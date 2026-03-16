//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 인라인 스키마의 properties에서 응답 필드를 수집
package validator

func collectInlineResponseFields(schema openAPISchema, schemas map[string]openAPISchema, fields map[string]FieldSymbol) {
	for fname, fprop := range schema.Properties {
		fields[fname] = toFieldSymbol(fprop, schemas)
	}
}
