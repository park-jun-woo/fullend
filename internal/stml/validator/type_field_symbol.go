//ff:type feature=stml-validate type=model
//ff:what 응답 필드의 타입 정보를 나타내는 심볼
package validator

// FieldSymbol represents a response field with type info.
type FieldSymbol struct {
	Type     string // "string", "integer", "array", "object"
	ItemType string // item type if array
}
