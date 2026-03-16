//ff:func feature=stml-validate type=util control=iteration dimension=1
//ff:what $ref를 해석하여 응답 필드 심볼을 수집
package validator

// resolveResponseFields resolves a $ref and collects response field symbols.
func resolveResponseFields(schemas map[string]openAPISchema, ref string, fields map[string]FieldSymbol) {
	name := refName(ref)
	schema, ok := schemas[name]
	if !ok {
		return
	}
	for fname, fprop := range schema.Properties {
		fields[fname] = toFieldSymbol(fprop, schemas)
	}
}
