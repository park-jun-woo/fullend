//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=http-handler
//ff:what 시퀀스에서 result 변수의 타입과 출처 정보를 수집
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectResultInfo(seqs []ssacparser.Sequence) (map[string]string, map[string]varSource) {
	resultTypes := map[string]string{}
	varSources := map[string]varSource{}
	for _, seq := range seqs {
		if seq.Result == nil {
			continue
		}
		resultTypes[seq.Result.Var] = seq.Result.Type
		varSources[seq.Result.Var] = resolveVarSource(seq)
	}
	return resultTypes, varSources
}
