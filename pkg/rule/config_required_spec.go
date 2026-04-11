//ff:type feature=rule type=model
//ff:what ConfigRequiredSpec — ConfigRequired 규칙의 판정 기준
package rule

// ConfigRequiredSpec configures a ConfigRequired rule.
// ConfigKey selects which Ground.Config key must be set.
type ConfigRequiredSpec struct {
	BaseSpec
	ConfigKey string
}
