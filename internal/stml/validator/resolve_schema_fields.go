//ff:func feature=stml-validate type=util control=iteration dimension=1
//ff:what $ref를 해석하여 요청 필드를 수집
package validator

// resolveSchemaFields resolves a $ref and collects field names into the map.
func resolveSchemaFields(schemas map[string]openAPISchema, ref string, fields map[string]string) {
	name := refName(ref)
	schema, ok := schemas[name]
	if !ok {
		return
	}
	for fname, fprop := range schema.Properties {
		typ := fprop.Type
		if fprop.Ref != "" {
			typ = "object"
		}
		fields[fname] = typ
	}
}
