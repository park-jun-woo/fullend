//ff:func feature=contract type=util control=sequence
//ff:what 디렉티브를 JS 코멘트 문자열로 변환한다
package contract

import "fmt"

// StringJS returns the directive as a JS comment (with space after //).
func (d *Directive) StringJS() string {
	return fmt.Sprintf("// fullend:%s ssot=%s contract=%s", d.Ownership, d.SSOT, d.Contract)
}
