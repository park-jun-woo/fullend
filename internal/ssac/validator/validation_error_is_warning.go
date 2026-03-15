//ff:method feature=ssac-validate type=model
//ff:what ValidationError.IsWarning() 메서드
package validator

// IsWarning은 이 에러가 경고인지 반환한다.
func (e ValidationError) IsWarning() bool {
	return e.Level == "WARNING"
}
