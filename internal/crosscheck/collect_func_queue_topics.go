//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=queue-check
//ff:what 단일 SSaC 함수에서 publish/subscribe 토픽을 수집
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func collectFuncQueueTopics(fn ssacparser.ServiceFunc, publishTopics map[string]map[string]bool, subscribeTopics map[string]ssacparser.ServiceFunc) {
	if fn.Subscribe != nil {
		subscribeTopics[fn.Subscribe.Topic] = fn
	}
	for _, seq := range fn.Sequences {
		if seq.Type == "publish" {
			publishTopics[seq.Topic] = collectInputKeys(seq.Inputs)
		}
	}
}
