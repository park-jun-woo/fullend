//ff:type feature=rule type=model
//ff:what TypeMatchSpec — TypeMatch 규칙의 판정 기준
package rule

// TypeMatchSpec configures a TypeMatch rule.
// LookupKey is the prefix into Ground.Types (appends ".name").
type TypeMatchSpec struct {
	BaseSpec
	LookupKey string
}
