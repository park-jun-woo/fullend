//ff:type feature=symbol type=model
//ff:what 모델 메서드의 상세 정보
package validator

// MethodInfo는 모델 메서드의 상세 정보다.
type MethodInfo struct {
	Cardinality string            // "one", "many", "exec"
	Params      []string          // interface 파라미터명 (context.Context 제외, 패키지 모델용)
	ParamTypes  map[string]string // 파라미터명 → Go 타입 (e.g. "amount" → "int"). @call Request struct 필드용
	ErrStatus   int               // @error 어노테이션 값 (0이면 미지정)
}
