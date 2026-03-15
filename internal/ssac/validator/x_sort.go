//ff:type feature=symbol type=model
//ff:what x-sort 확장
package validator

// XSort는 x-sort 확장이다.
type XSort struct {
	Allowed   []string `yaml:"allowed"`
	Default   string   `yaml:"default"`
	Direction string   `yaml:"direction"`
}
