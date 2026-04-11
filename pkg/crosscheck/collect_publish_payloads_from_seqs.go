//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectPublishPayloadsFromSeqs — 시퀀스 목록에서 @publish topic별 필드 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectPublishPayloadsFromSeqs(payloads map[string][]string, seqs []ssac.Sequence) {
	for _, seq := range seqs {
		if seq.Type == "publish" && seq.Topic != "" {
			payloads[seq.Topic] = collectFieldKeys(seq.Fields)
		}
	}
}
