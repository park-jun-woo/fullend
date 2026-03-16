//ff:func feature=stml-validate type=rule control=sequence
//ff:what 컴포넌트 파일이 없을 때 오류 생성
package validator

import "fmt"

func errComponentNotFound(file, comp, path string) ValidationError {
	return ValidationError{file, fmt.Sprintf("data-component=%q", comp), fmt.Sprintf("%s 파일이 없습니다", path)}
}
