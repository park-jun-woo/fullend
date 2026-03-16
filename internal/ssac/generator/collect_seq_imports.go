//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 시퀀스에서 필요한 패키지 import을 수집
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func collectSeqImports(sf parser.ServiceFunc, seen map[string]bool) {
	for _, seq := range sf.Sequences {
		if seq.Type == parser.SeqState {
			seen["states/"+seq.DiagramID+"state"] = true
		}
		if seq.Type == parser.SeqAuth {
			seen["authz"] = true
		}
		if seq.Type == parser.SeqPublish {
			seen["queue"] = true
		}
		if seq.Result != nil && seq.Result.Wrapper != "" && !hasDirectResponse(sf.Sequences) {
			seen["github.com/geul-org/fullend/pkg/pagination"] = true
		}
	}
}
