//ff:type feature=symbol type=model
//ff:what openAPIMediaType 타입 정의
package validator

type openAPIMediaType struct {
	Schema openAPISchema `yaml:"schema"`
}
