//ff:type feature=symbol type=model
//ff:what x-filter 확장
package validator

// XFilter는 x-filter 확장이다.
type XFilter struct {
	Allowed []string `yaml:"allowed"`
}
