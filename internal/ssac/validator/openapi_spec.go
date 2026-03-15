//ff:type feature=symbol type=model
//ff:what openAPISpec 타입 정의
package validator

type openAPISpec struct {
	Paths      map[string]openAPIPathItem `yaml:"paths"`
	Components openAPIComponents          `yaml:"components"`
}
