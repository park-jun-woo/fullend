//ff:func feature=stml-validate type=rule control=sequence
//ff:what 응답 스키마와 custom.ts에 바인딩 필드가 없을 때 오류 생성
package validator

import "fmt"

func errBindNotFound(file, op, field string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-bind=%q", field), fmt.Sprintf("%q의 response schema에도, custom.ts에도 %q가 없습니다", op, field)}
}
