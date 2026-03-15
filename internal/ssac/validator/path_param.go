//ff:type feature=symbol type=model
//ff:what OpenAPI path parameter
package validator

// PathParam은 OpenAPI path parameter다.
type PathParam struct {
	Name   string // 원본 이름 (e.g. "CourseID")
	GoType string // Go 타입 (e.g. "int64")
}
