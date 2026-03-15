//ff:type feature=symbol type=model
//ff:what openAPIComponents 타입 정의
package validator

type openAPIComponents struct {
	Schemas map[string]openAPISchema `yaml:"schemas"`
}
