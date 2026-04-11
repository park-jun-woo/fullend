//ff:type feature=manifest type=model
//ff:what FieldConstraint — OpenAPI 스키마 property의 제약조건
package openapi

// FieldConstraint holds constraints for a single OpenAPI schema property.
type FieldConstraint struct {
	Type      string
	Format    string
	MaxLength *int
	MinLength *int
	Enum      []string
	Required  bool
}
