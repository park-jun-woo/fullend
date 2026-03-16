//ff:type feature=contract type=model
//ff:what 계약 해시 불일치 경고를 나타내는 구조체
package contract

// Warning records a contract mismatch between preserved body and regenerated contract.
type Warning struct {
	File        string
	Function    string
	OldContract string
	NewContract string
}
