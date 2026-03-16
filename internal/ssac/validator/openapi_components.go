//ff:type feature=symbol type=model topic=openapi
//ff:what openAPIComponents 타입 정의
package validator

type openAPIComponents struct {
	Schemas map[string]openAPISchema `yaml:"schemas"`
}
