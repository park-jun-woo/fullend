//ff:type feature=stml-validate type=model
//ff:what OpenAPI 응답 YAML 구조체
package validator

type openAPIResponse struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}
