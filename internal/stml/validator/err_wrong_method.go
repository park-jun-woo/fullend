//ff:func feature=stml-validate type=rule control=sequence
//ff:what HTTP 메서드 불일치 오류 생성
package validator

import "fmt"

func errWrongMethod(file, attr, op, got, want string) ValidationError {
	return ValidationError{file, attr, fmt.Sprintf("%q은 %s 메서드입니다 (%s이어야 함)", op, got, want)}
}
