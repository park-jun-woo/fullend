//ff:type feature=ssac-parse type=model
//ff:what 함수 호출 인자 타입
package ssac

// Arg는 함수 호출 인자다.
type Arg struct {
	Source  string // "request", 변수명, 또는 "" (리터럴)
	Field   string // "CourseID", "ID" 등
	Literal string // "cancelled" 등 (Source가 ""일 때)
}
