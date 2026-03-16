//ff:type feature=stml-validate type=model
//ff:what OpenAPI 파라미터 YAML 구조체
package validator

type openAPIParam struct {
	Name   string        `yaml:"name"`
	In     string        `yaml:"in"`
	Schema openAPISchema `yaml:"schema"`
}
