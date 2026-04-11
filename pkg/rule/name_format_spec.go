//ff:type feature=rule type=model
//ff:what NameFormatSpec — NameFormat 규칙의 판정 기준
package rule

// NameFormatSpec configures a NameFormat rule.
// Pattern: "uppercase-start", "no-dot-prefix", "dot-method".
type NameFormatSpec struct {
	BaseSpec
	Pattern string
}
