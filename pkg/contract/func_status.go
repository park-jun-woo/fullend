//ff:type feature=contract type=model
//ff:what 함수별 계약 상태를 나타내는 구조체
package contract

// FuncStatus describes the contract status of a single function.
type FuncStatus struct {
	File      string // relative path from artifacts dir
	Function  string
	Directive Directive
	Status    string // "gen", "preserve", "broken", "orphan"
	Detail    string // violation detail
}
