//ff:type feature=ssac-validate type=model
//ff:what ValidationError 타입 정의
package validator

// ValidationError는 검증 에러 하나를 나타낸다.
type ValidationError struct {
	FileName string // 원본 파일명
	FuncName string // 함수명
	SeqIndex int    // sequence 인덱스
	Tag      string // 관련 태그 (e.g. "@model", "@action")
	Message  string // 에러 메시지
	Level    string // "ERROR" 또는 "WARNING" (빈 문자열이면 ERROR)
}
