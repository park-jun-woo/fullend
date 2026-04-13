//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=model-collect
//ff:what 단일 ServiceFunc에서 모델 사용 정보를 수집
package ssac

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func collectModelUsagesFromFunc(sf ssacparser.ServiceFunc) []modelUsage {
	var usages []modelUsage
	for _, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == ssacparser.SeqCall {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 {
			continue
		}
		usages = append(usages, modelUsage{
			ModelName:  parts[0],
			MethodName: parts[1],
			Inputs:     seq.Inputs,
			Result:     seq.Result,
			FuncName:   sf.Name,
		})
	}
	return usages
}
