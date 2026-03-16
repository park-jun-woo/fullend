//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=queue-check
//ff:what SSaC 함수에서 publish/subscribe 토픽을 수집
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func collectQueueTopics(funcs []ssacparser.ServiceFunc) (map[string]map[string]bool, map[string]ssacparser.ServiceFunc) {
	publishTopics := map[string]map[string]bool{}
	subscribeTopics := map[string]ssacparser.ServiceFunc{}

	for _, fn := range funcs {
		collectFuncQueueTopics(fn, publishTopics, subscribeTopics)
	}

	return publishTopics, subscribeTopics
}
