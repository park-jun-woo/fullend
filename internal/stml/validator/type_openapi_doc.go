//ff:type feature=stml-validate type=model
//ff:what OpenAPI YAML 문서 최상위 구조체
package validator

type openAPIDoc struct {
	Paths      map[string]map[string]openAPIOperation `yaml:"paths"`
	Components struct {
		Schemas map[string]openAPISchema `yaml:"schemas"`
	} `yaml:"components"`
}
