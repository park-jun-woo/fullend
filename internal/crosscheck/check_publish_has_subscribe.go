//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=queue-check
//ff:what publish 토픽에 대한 subscribe 함수 존재 여부 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func checkPublishHasSubscribe(publishTopics map[string]map[string]bool, subscribeTopics map[string]ssacparser.ServiceFunc) []CrossError {
	var errs []CrossError
	for topic := range publishTopics {
		if _, ok := subscribeTopics[topic]; !ok {
			errs = append(errs, CrossError{
				Rule:       "Queue publish → subscribe",
				Context:    topic,
				Message:    fmt.Sprintf("@publish 토픽 %q에 대한 @subscribe 함수가 없습니다", topic),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("토픽 %q를 구독하는 @subscribe 함수를 추가하세요", topic),
			})
		}
	}
	return errs
}
