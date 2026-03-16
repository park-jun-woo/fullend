//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what subscribe 메시지 필드가 publish payload에 존재하는지 검증
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func checkQueueFieldMismatch(subscribeTopics map[string]ssacparser.ServiceFunc, publishTopics map[string]map[string]bool) []CrossError {
	var errs []CrossError
	for topic, fn := range subscribeTopics {
		pubFields, ok := publishTopics[topic]
		if !ok {
			continue
		}
		errs = append(errs, checkSubscribeFields(fn, topic, pubFields)...)
	}
	return errs
}
