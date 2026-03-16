//ff:type feature=symbol type=model topic=openapi
//ff:what openAPISchema 타입 정의
package validator

type openAPISchema struct {
	Type       string                   `yaml:"type"`
	Format     string                   `yaml:"format"`
	Properties map[string]openAPISchema `yaml:"properties"`
	Ref        string                   `yaml:"$ref"`
	Required   []string                 `yaml:"required"`
	MinLength  *int                     `yaml:"minLength"`
	MaxLength  *int                     `yaml:"maxLength"`
	Minimum    *float64                 `yaml:"minimum"`
	Maximum    *float64                 `yaml:"maximum"`
	Pattern    string                   `yaml:"pattern"`
	Enum       []interface{}            `yaml:"enum"`
}
