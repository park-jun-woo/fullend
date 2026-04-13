//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=import-collect
//ff:what 시퀀스에서 필요한 패키지 import을 수집
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectSeqImports(sf ssacparser.ServiceFunc, seen map[string]bool) {
	for _, seq := range sf.Sequences {
		if seq.Type == ssacparser.SeqState {
			seen["states/"+seq.DiagramID+"state"] = true
		}
		if seq.Type == ssacparser.SeqAuth {
			seen["authz"] = true
		}
		if seq.Type == ssacparser.SeqPublish {
			seen["queue"] = true
		}
		if seq.Result != nil && seq.Result.Wrapper != "" && !hasDirectResponse(sf.Sequences) {
			seen["github.com/park-jun-woo/fullend/pkg/pagination"] = true
		}
	}
}
