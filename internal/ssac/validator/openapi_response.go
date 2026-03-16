//ff:type feature=symbol type=model topic=openapi
//ff:what openAPIResponse 타입 정의
package validator

type openAPIResponse struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}
