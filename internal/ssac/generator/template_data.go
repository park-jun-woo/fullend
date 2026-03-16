//ff:type feature=ssac-gen type=model topic=template-data
//ff:what 템플릿 실행에 필요한 시퀀스별 데이터 구조체
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

type templateData struct {
	// 공통
	Message  string
	FirstErr bool

	// get/post/put/delete
	ModelCall string // "courseModel.FindByID"
	ArgsCode  string // "courseID, currentUser.ID"
	Result    *parser.Result

	// empty/exists
	Target      string
	ZeroCheck   string
	ExistsCheck string

	// state
	DiagramID   string
	Transition  string
	InputFields string // "Status: reservation.Status, ..."

	// auth
	Action   string
	Resource string

	// call
	PkgName    string
	FuncMethod string
	ErrStatus  string // "http.StatusInternalServerError", "http.StatusUnauthorized" 등

	// publish
	Topic      string // "order.completed"
	OptionCode string // ", queue.WithDelay(1800)" 또는 ""

	// response
	ResponseFields map[string]string

	// list
	HasTotal bool

	// reassign: result var already declared -> use = instead of :=
	ReAssign bool

	// unused: result var not referenced later -> use _ instead of var name
	Unused bool

	// errDeclared: err variable already declared before this sequence
	ErrDeclared bool
}
