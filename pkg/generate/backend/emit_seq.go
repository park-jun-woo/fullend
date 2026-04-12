//ff:func feature=rule type=generator control=selection
//ff:what emitSeq — Trace 패턴으로 시퀀스 타입별 코드 생성기 디스패치
package backend

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/generate/trace"
)

func emitSeq(seq parsessac.Sequence, p trace.Pattern) string {
	switch {
	case p.Has("HasPaginateCursor"):
		return emitPaginateCursor(seq)
	case p.Has("HasPaginateOffset"):
		return emitPaginateOffset(seq)
	case p.Has("HasSliceResult"):
		return emitSimpleGet(seq)
	case p.Has("HasFKRef") && p.Has("IsGet"):
		return emitFKGet(seq)
	case p.Has("IsGet"):
		return emitSimpleGet(seq)
	case p.Has("IsPost"):
		return emitPost(seq)
	case p.Has("IsPut"):
		return emitPut(seq)
	case p.Has("IsDelete"):
		return emitDelete(seq)
	case p.Has("IsEmpty"):
		return emitEmpty(seq)
	case p.Has("IsExists"):
		return emitExists(seq)
	case p.Has("IsState"):
		return emitState(seq)
	case p.Has("IsAuth"):
		return emitAuth(seq)
	case p.Has("IsCall"):
		return emitCall(seq)
	case p.Has("IsPublish"):
		return emitPublish(seq)
	case p.Has("IsResponse"):
		return emitResponse(seq)
	}
	return ""
}
