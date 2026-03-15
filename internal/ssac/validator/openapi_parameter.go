//ff:type feature=symbol type=model
//ff:what openAPIParameter 타입 정의
package validator

type openAPIParameter struct {
	Name   string        `yaml:"name"`
	In     string        `yaml:"in"`
	Schema openAPISchema `yaml:"schema"`
}
