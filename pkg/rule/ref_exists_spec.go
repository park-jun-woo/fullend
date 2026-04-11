//ff:type feature=rule type=model
//ff:what RefExistsSpec — RefExists 규칙의 판정 기준
package rule

// RefExistsSpec configures a RefExists rule.
// LookupKey selects which Ground.Lookup set to check against.
type RefExistsSpec struct {
	BaseSpec
	LookupKey string
}
