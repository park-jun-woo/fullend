//ff:func feature=stml-validate type=rule control=sequence
//ff:what data-each 필드가 응답에 없을 때 오류 생성
package validator

import "fmt"

func errEachNotFound(file, op, field string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-each=%q", field), fmt.Sprintf("%q의 response에 %q 필드가 없습니다", op, field)}
}
