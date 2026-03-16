//ff:func feature=contract type=util control=sequence
//ff:what 디렉티브를 Go 코멘트 문자열로 변환한다
package contract

import "fmt"

// String returns the directive as a Go comment.
func (d *Directive) String() string {
	return fmt.Sprintf("//fullend:%s ssot=%s contract=%s", d.Ownership, d.SSOT, d.Contract)
}
