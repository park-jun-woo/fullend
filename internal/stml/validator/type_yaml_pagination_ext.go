//ff:type feature=stml-validate type=model
//ff:what x-pagination YAML 확장 구조체
package validator

type yamlPaginationExt struct {
	Style        string `yaml:"style"`
	DefaultLimit int    `yaml:"defaultLimit"`
	MaxLimit     int    `yaml:"maxLimit"`
}
