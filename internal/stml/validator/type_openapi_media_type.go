//ff:type feature=stml-validate type=model
//ff:what OpenAPI 미디어 타입 YAML 구조체
package validator

type openAPIMediaType struct {
	Schema openAPISchema `yaml:"schema"`
}
