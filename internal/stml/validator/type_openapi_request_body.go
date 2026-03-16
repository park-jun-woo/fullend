//ff:type feature=stml-validate type=model
//ff:what OpenAPI 요청 본문 YAML 구조체
package validator

type openAPIRequestBody struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}
