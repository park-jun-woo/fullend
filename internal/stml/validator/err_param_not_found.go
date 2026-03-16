//ff:func feature=stml-validate type=rule control=sequence
//ff:what 파라미터가 OpenAPI에 없을 때 오류 생성
package validator

import "fmt"

func errParamNotFound(file, op, param string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-param-%s", param), fmt.Sprintf("%q의 parameters에 %q가 없습니다", op, param)}
}
