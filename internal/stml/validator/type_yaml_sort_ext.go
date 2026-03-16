//ff:type feature=stml-validate type=model
//ff:what x-sort YAML 확장 구조체
package validator

type yamlSortExt struct {
	Allowed   []string `yaml:"allowed"`
	Default   string   `yaml:"default"`
	Direction string   `yaml:"direction"`
}
