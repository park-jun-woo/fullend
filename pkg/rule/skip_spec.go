//ff:type feature=rule type=model
//ff:what SkipSpec — IsSkipped defeater의 판정 기준
package rule

// SkipSpec configures an IsSkipped defeater.
// Kind is the SSOT kind to check (e.g., "DDL", "SSaC").
type SkipSpec struct {
	BaseSpec
	Kind string
}
