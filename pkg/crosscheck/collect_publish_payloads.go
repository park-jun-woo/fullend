//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectPublishPayloads — @publish 시퀀스에서 topic별 payload 필드 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func collectPublishPayloads(fs *fullend.Fullstack) map[string][]string {
	payloads := make(map[string][]string)
	for _, fn := range fs.ServiceFuncs {
		collectPublishPayloadsFromSeqs(payloads, fn.Sequences)
	}
	return payloads
}
