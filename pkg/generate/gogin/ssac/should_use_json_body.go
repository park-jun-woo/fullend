//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=request-params
//ff:what OpenAPI requestBody 또는 @post/@put 시퀀스로 JSON body 사용 여부를 판단
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func shouldUseJSONBody(seqs []ssacparser.Sequence, st *rule.Ground, operationID string, rawParams []rawParam) bool {
	if len(rawParams) == 0 {
		return false
	}
	// 1차: OpenAPI에 requestBody가 있으면 JSON body
	if st != nil {
		if op, ok := st.Ops[operationID]; ok {
			return op.HasRequestBody
		}
	}
	// 2차 fallback: Operations 미등록 시 @post/@put 시퀀스로 판단 (테스트 호환)
	for _, seq := range seqs {
		if seq.Type == ssacparser.SeqPost || seq.Type == ssacparser.SeqPut {
			return true
		}
	}
	return false
}
