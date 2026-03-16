//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 단일 ServiceFunc에서 사용하는 모델명을 도메인 맵에 추가
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

func collectModelsForFunc(sf parser.ServiceFunc, domainSet map[string]map[string]bool) {
	domain := sf.Domain
	if domainSet[domain] == nil {
		domainSet[domain] = map[string]bool{}
	}
	for _, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == parser.SeqCall {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) >= 1 {
			domainSet[domain][parts[0]] = true
		}
	}
}
