//ff:type feature=policy type=model
//ff:what @ownership 어노테이션에서 추출한 소유권 매핑 구조체
package policy

// OwnershipMapping represents a @ownership annotation.
type OwnershipMapping struct {
	Resource  string
	Table     string
	Column    string
	JoinTable string // empty if direct lookup
	JoinFK    string // empty if direct lookup
}
