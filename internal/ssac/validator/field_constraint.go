//ff:type feature=symbol type=model topic=openapi
//ff:what OpenAPI schema property의 검증 제약
package validator

// FieldConstraint는 OpenAPI schema property의 검증 제약을 담는다.
type FieldConstraint struct {
	Required  bool
	Format    string
	MinLength *int
	MaxLength *int
	Minimum   *float64
	Maximum   *float64
	Pattern   string
	Enum      []string
}
