//ff:type feature=symbol type=model
//ff:what x-include 확장
package validator

// XInclude는 x-include 확장이다.
type XInclude struct {
	Allowed []string `yaml:"allowed"`
}
