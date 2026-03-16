//ff:func feature=symbol type=util control=iteration dimension=2 topic=openapi
//ff:what requestBody 스키마에서 필드별 검증 제약을 수집
package validator

func extractRequestSchema(schema openAPISchema, allSchemas map[string]openAPISchema) RequestSchema {
	resolved := resolveSchema(schema, allSchemas)
	requiredSet := map[string]bool{}
	for _, r := range resolved.Required {
		requiredSet[r] = true
	}
	rs := RequestSchema{Fields: map[string]FieldConstraint{}}
	for name, prop := range resolved.Properties {
		prop = resolveSchema(prop, allSchemas)
		fc := FieldConstraint{
			Required:  requiredSet[name],
			Format:    prop.Format,
			MinLength: prop.MinLength,
			MaxLength: prop.MaxLength,
			Minimum:   prop.Minimum,
			Maximum:   prop.Maximum,
			Pattern:   prop.Pattern,
		}
		for _, e := range prop.Enum {
			if s, ok := e.(string); ok {
				fc.Enum = append(fc.Enum, s)
			}
		}
		rs.Fields[name] = fc
	}
	return rs
}
