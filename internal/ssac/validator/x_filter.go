//ff:type feature=symbol type=model topic=openapi
//ff:what x-filter 확장
package validator

// XFilter는 x-filter 확장이다.
type XFilter struct {
	Allowed []string `yaml:"allowed"`
}
