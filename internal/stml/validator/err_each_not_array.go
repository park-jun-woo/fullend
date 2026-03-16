//ff:func feature=stml-validate type=rule control=sequence
//ff:what data-each 필드가 배열이 아닐 때 오류 생성
package validator

import "fmt"

func errEachNotArray(file, op, field string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-each=%q", field), fmt.Sprintf("%q의 response에서 %q는 배열이 아닙니다", op, field)}
}
