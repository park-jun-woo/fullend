//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=http-handler
//ff:what 시퀀스에서 참조되는 변수명을 수집
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// collectUsedVars는 시퀀스에서 참조되는 변수명을 수집한다.
func collectUsedVars(seqs []ssacparser.Sequence) map[string]bool {
	used := map[string]bool{}
	for _, seq := range seqs {
		collectUsedVarsFromSeq(seq, used)
	}
	return used
}
