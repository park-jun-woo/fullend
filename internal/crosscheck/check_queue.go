//ff:func feature=crosscheck type=rule control=sequence topic=queue-check
//ff:what publish ↔ subscribe 토픽 교차 참조와 설정 유무를 검증
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckQueue validates publish ↔ subscribe cross-references.
func CheckQueue(funcs []ssacparser.ServiceFunc, queueBackend string) []CrossError {
	var errs []CrossError

	publishTopics, subscribeTopics := collectQueueTopics(funcs)

	hasQueue := len(publishTopics) > 0 || len(subscribeTopics) > 0

	if hasQueue && queueBackend == "" {
		errs = append(errs, CrossError{
			Rule:       "Queue config",
			Context:    "fullend.yaml",
			Message:    "fullend.yaml에 queue 설정이 없지만 @publish/@subscribe가 사용되었습니다",
			Level:      "ERROR",
			Suggestion: "fullend.yaml에 queue.backend을 설정하세요 (예: postgres, memory)",
		})
	}

	errs = append(errs, checkPublishHasSubscribe(publishTopics, subscribeTopics)...)
	errs = append(errs, checkSubscribeHasPublish(subscribeTopics, publishTopics)...)
	errs = append(errs, checkQueueFieldMismatch(subscribeTopics, publishTopics)...)

	return errs
}
