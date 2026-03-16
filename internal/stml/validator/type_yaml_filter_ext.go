//ff:type feature=stml-validate type=model
//ff:what x-filter YAML 확장 구조체
package validator

type yamlFilterExt struct {
	Allowed []string `yaml:"allowed"`
}
