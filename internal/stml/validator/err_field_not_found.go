//ff:func feature=stml-validate type=rule control=sequence
//ff:what 요청 스키마에 필드가 없을 때 오류 생성
package validator

import "fmt"

func errFieldNotFound(file, op, field string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-field=%q", field), fmt.Sprintf("%q의 request schema에 %q 필드가 없습니다", op, field)}
}
