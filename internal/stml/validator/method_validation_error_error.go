//ff:func feature=stml-validate type=model control=sequence
//ff:what ValidationError를 문자열로 변환하는 메서드
package validator

import "fmt"

func (e ValidationError) Error() string {
	return fmt.Sprintf("ERROR: %s — %s: %s", e.File, e.Attr, e.Message)
}
