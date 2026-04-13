//ff:type feature=contract type=model
//ff:what 보존된 함수의 디렉티브와 본문을 담는 구조체
package contract

// PreservedFunc holds a preserved function's directive and body text.
type PreservedFunc struct {
	Directive Directive
	BodyText  string // raw source between { and }, excluding the braces themselves
}
