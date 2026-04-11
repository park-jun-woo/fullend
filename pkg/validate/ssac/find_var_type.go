//ff:func feature=rule type=util control=iteration dimension=1
//ff:what findVarType — 시퀀스 목록에서 변수명으로 result 타입 조회
package ssac

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func findVarType(seqs []parsessac.Sequence, varName string) string {
	for _, seq := range seqs {
		if seq.Result != nil && seq.Result.Var == varName {
			return seq.Result.Type
		}
	}
	return ""
}
