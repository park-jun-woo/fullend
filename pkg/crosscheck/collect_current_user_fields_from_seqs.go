//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectCurrentUserFieldsFromSeqs — 시퀀스 목록에서 currentUser 필드 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectCurrentUserFieldsFromSeqs(seqs []ssac.Sequence, seen map[string]bool) {
	for _, seq := range seqs {
		collectCurrentUserFieldsFromArgs(seq.Args, seen)
	}
}
