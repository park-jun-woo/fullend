//ff:type feature=rule type=model
//ff:what ModelRefExistsSpec — ModelRefExists 규칙의 판정 기준
package rule

// ModelRefExistsSpec configures a ModelRefExists rule.
// Ground.Models 를 조회하므로 LookupKey 가 불필요하다.
type ModelRefExistsSpec struct {
	BaseSpec
}
