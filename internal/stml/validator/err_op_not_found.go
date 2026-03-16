//ff:func feature=stml-validate type=rule control=sequence
//ff:what operationId가 OpenAPI에 없을 때 오류 생성
package validator

import "fmt"

func errOpNotFound(file, attr, op string) ValidationError {
	return ValidationError{file, attr, fmt.Sprintf("OpenAPI에 %q operationId가 없습니다", op)}
}
