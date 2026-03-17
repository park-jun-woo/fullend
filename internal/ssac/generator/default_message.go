//ff:func feature=ssac-gen type=util control=selection topic=template-data
//ff:what 시퀀스 타입별 기본 에러 메시지를 반환
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func defaultMessage(seq parser.Sequence) string {
	modelName := ""
	if seq.Model != "" {
		parts := strings.SplitN(seq.Model, ".", 2)
		modelName = parts[0]
	}

	switch seq.Type {
	case parser.SeqGet:
		return modelName + " 조회 실패"
	case parser.SeqPost:
		return modelName + " 생성 실패"
	case parser.SeqPut:
		return modelName + " 수정 실패"
	case parser.SeqDelete:
		return modelName + " 삭제 실패"
	case parser.SeqEmpty:
		return seq.Target + "가 존재하지 않습니다"
	case parser.SeqExists:
		return seq.Target + "가 이미 존재합니다"
	case parser.SeqState:
		return "상태 전이가 허용되지 않습니다"
	case parser.SeqAuth:
		return "권한이 없습니다"
	case parser.SeqCall:
		return "호출 실패"
	}
	return "처리 실패"
}
