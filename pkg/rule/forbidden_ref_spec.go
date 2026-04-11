//ff:type feature=rule type=model
//ff:what ForbiddenRefSpec — ForbiddenRef 규칙의 판정 기준
package rule

// ForbiddenRefSpec configures a ForbiddenRef rule.
// LookupKey selects which Ground.Lookup set holds forbidden names.
type ForbiddenRefSpec struct {
	BaseSpec
	LookupKey string
}
