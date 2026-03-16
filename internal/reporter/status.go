//ff:type feature=reporter type=model
//ff:what 검증 단계의 결과 상태를 나타내는 타입과 상수
package reporter

// Status represents the outcome of a validation step.
type Status int

const (
	Pass Status = iota
	Fail
	Skip
)
