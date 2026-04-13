//ff:func feature=ssac-gen type=util control=selection topic=template-data
//ff:what 시퀀스 타입별 기본 에러 메시지를 반환
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func defaultMessage(seq ssacparser.Sequence) string {
	modelName := ""
	if seq.Model != "" {
		parts := strings.SplitN(seq.Model, ".", 2)
		modelName = parts[0]
	}

	switch seq.Type {
	case ssacparser.SeqGet:
		return modelName + " 조회 실패"
	case ssacparser.SeqPost:
		return modelName + " 생성 실패"
	case ssacparser.SeqPut:
		return modelName + " 수정 실패"
	case ssacparser.SeqDelete:
		return modelName + " 삭제 실패"
	case ssacparser.SeqEmpty:
		return seq.Target + "가 존재하지 않습니다"
	case ssacparser.SeqExists:
		return seq.Target + "가 이미 존재합니다"
	case ssacparser.SeqState:
		return "상태 전이가 허용되지 않습니다"
	case ssacparser.SeqAuth:
		return "권한이 없습니다"
	case ssacparser.SeqCall:
		return "호출 실패"
	}
	return "처리 실패"
}
