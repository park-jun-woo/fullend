//ff:type feature=rule type=model
//ff:what SeqClaim — 시퀀스 코드 생성 시 Toulmin Context로 전달할 claim
package backend

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// SeqClaim bundles a sequence with its surrounding function context.
type SeqClaim struct {
	Type       string
	Seq        parsessac.Sequence
	FuncName   string
	HasQuery   bool
	FKRef      bool
	IsBuiltin  bool
}
