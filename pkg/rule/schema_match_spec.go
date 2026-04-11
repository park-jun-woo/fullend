//ff:type feature=rule type=model
//ff:what SchemaMatchSpec — SchemaMatch 규칙의 판정 기준
package rule

// SchemaMatchSpec configures a SchemaMatch rule.
// LookupKey selects which Ground.Schemas to check against.
type SchemaMatchSpec struct {
	BaseSpec
	LookupKey string
}
