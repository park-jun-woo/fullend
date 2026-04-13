//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=model-collect
//ff:what 단일 ServiceFunc에서 사용하는 모델명을 도메인 맵에 추가
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func collectModelsForFunc(sf ssacparser.ServiceFunc, domainSet map[string]map[string]bool) {
	domain := sf.Feature
	if domainSet[domain] == nil {
		domainSet[domain] = map[string]bool{}
	}
	for _, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == ssacparser.SeqCall {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) >= 1 {
			domainSet[domain][parts[0]] = true
		}
	}
}
