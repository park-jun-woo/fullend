//ff:func feature=stml-validate type=rule control=sequence
//ff:what x-pagination 미선언 시 오류 생성
package validator

import "fmt"

func errPaginateNoExt(file, op string) ValidationError {
	return ValidationError{file, "data-paginate", fmt.Sprintf("%q에 x-pagination이 선언되지 않았습니다", op)}
}
