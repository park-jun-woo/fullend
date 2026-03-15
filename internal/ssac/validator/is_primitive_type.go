//ff:func feature=ssac-validate type=util
//ff:what Go 기본 타입 여부를 반환한다
package validator

// primitiveTypes는 Go 기본 타입 집합이다.
var primitiveTypes = map[string]bool{
	"string": true, "int": true, "int32": true, "int64": true,
	"float32": true, "float64": true, "bool": true, "byte": true,
	"rune": true, "time.Time": true,
}

// isPrimitiveType는 Go 기본 타입 여부를 반환한다.
func isPrimitiveType(s string) bool {
	return primitiveTypes[s]
}
