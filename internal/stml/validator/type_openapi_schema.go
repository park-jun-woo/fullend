//ff:type feature=stml-validate type=model
//ff:what OpenAPI 스키마 YAML 구조체
package validator

type openAPISchema struct {
	Ref        string                   `yaml:"$ref"`
	Type       string                   `yaml:"type"`
	Properties map[string]openAPISchema `yaml:"properties"`
	Items      *openAPISchema           `yaml:"items"`
}
