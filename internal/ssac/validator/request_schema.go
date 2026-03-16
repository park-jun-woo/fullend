//ff:type feature=symbol type=model topic=openapi
//ff:what operationId별 requestBody 필드 제약
package validator

// RequestSchema는 하나의 operationId에 대한 requestBody 필드별 제약을 담는다.
type RequestSchema struct {
	Fields map[string]FieldConstraint
}
