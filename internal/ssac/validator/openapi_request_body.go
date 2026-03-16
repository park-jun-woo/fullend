//ff:type feature=symbol type=model topic=openapi
//ff:what openAPIRequestBody 타입 정의
package validator

type openAPIRequestBody struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}
