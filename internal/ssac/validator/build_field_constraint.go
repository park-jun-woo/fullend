//ff:func feature=symbol type=util control=iteration dimension=1 topic=openapi
//ff:what 단일 OpenAPI 프로퍼티에서 FieldConstraint를 생성한다

package validator

// buildFieldConstraint creates a FieldConstraint from a resolved OpenAPI property.
func buildFieldConstraint(name string, prop openAPISchema, allSchemas map[string]openAPISchema, requiredSet map[string]bool) FieldConstraint {
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
	return fc
}
