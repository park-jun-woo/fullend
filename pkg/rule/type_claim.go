//ff:type feature=rule type=model
//ff:what TypeClaim — 타입 비교를 위한 이름+소스타입 쌍
package rule

// TypeClaim holds a name and its source type for TypeMatch comparison.
type TypeClaim struct {
	Name       string
	SourceType string
}
