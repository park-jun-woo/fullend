//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what subscribe 토픽에 대한 publish 존재 여부 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func checkSubscribeHasPublish(subscribeTopics map[string]ssacparser.ServiceFunc, publishTopics map[string]map[string]bool) []CrossError {
	var errs []CrossError
	for topic := range subscribeTopics {
		if _, ok := publishTopics[topic]; !ok {
			errs = append(errs, CrossError{
				Rule:       "Queue subscribe → publish",
				Context:    topic,
				Message:    fmt.Sprintf("@subscribe 토픽 %q에 대한 @publish가 없습니다", topic),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("토픽 %q를 발행하는 @publish 시퀀스를 추가하세요", topic),
			})
		}
	}
	return errs
}
