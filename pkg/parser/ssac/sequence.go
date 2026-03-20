//ff:type feature=ssac-parse type=model
//ff:what 하나의 시퀀스 라인 타입
package parser

// Sequence는 하나의 시퀀스 라인이다.
type Sequence struct {
	Type string // "get", "post", "put", "delete", "empty", "exists", "state", "auth", "call", "response"

	// get/post/put/delete/call 공통: 함수 호출
	Package string // "session" (패키지 접두사, 없으면 "")
	Model   string // "Course.FindByID" 또는 "auth.VerifyPassword"
	Args    []Arg  // 호출 인자

	// get/post/call: 대입
	Result *Result // 결과 바인딩 (nil이면 대입 없음)

	// empty/exists: guard
	Target string // "course" 또는 "course.InstructorID"

	// state: 상태 전이
	DiagramID  string            // "reservation"
	Inputs     map[string]string // {status: "reservation.Status"}
	Transition string            // "cancel"

	// publish: 이벤트 발행
	Topic   string            // "order.completed"
	Options map[string]string // {delay: "1800"} (선택)
	// Inputs 재사용: payload

	// auth: 권한 검사
	Action   string // "delete"
	Resource string // "project"
	// Inputs 재사용     // {id: "project.ID", owner: "project.OwnerID"}

	// response: 필드 매핑
	Fields map[string]string // {course: "course", instructor_name: "instructor.Name"}

	// 공통
	Message      string // 에러 메시지
	ErrStatus    int    // 에러 HTTP 상태 코드 (0이면 타입별 기본값: @call→500, @empty→404, @exists→409, @state→409, @auth→403)
	SuppressWarn bool   // @type! — WARNING 억제
}
