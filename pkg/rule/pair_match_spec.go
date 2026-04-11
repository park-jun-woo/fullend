//ff:type feature=rule type=model
//ff:what PairMatchSpec — PairMatch 규칙의 판정 기준
package rule

// PairMatchSpec configures a PairMatch rule.
// LookupKey selects which Ground.Pairs set to check against.
type PairMatchSpec struct {
	BaseSpec
	LookupKey string
}
