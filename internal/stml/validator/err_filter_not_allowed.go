//ff:func feature=stml-validate type=rule control=sequence
//ff:what 필터 컬럼이 허용 목록에 없을 때 오류 생성
package validator

import "fmt"

func errFilterNotAllowed(file, op, col string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-filter=%q", col), fmt.Sprintf("%q의 x-filter.allowed에 %q가 없습니다", op, col)}
}
