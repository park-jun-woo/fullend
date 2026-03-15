//ff:type feature=symbol type=model
//ff:what openAPISchema 타입 정의
package validator

type openAPISchema struct {
	Type       string                   `yaml:"type"`
	Format     string                   `yaml:"format"`
	Properties map[string]openAPISchema `yaml:"properties"`
	Ref        string                   `yaml:"$ref"`
}
