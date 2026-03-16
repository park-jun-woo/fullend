//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what subscribe 메시지 struct 필드가 publish payload 키에 존재하는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func checkSubscribeFields(fn ssacparser.ServiceFunc, topic string, pubFields map[string]bool) []CrossError {
	typeName := ""
	if fn.Param != nil {
		typeName = fn.Param.TypeName
	}
	if typeName == "" {
		return nil
	}

	structFields := findStructFields(fn.Structs, typeName)

	var errs []CrossError
	for _, fieldName := range structFields {
		if !pubFields[fieldName] {
			errs = append(errs, CrossError{
				Rule:    "Queue field mismatch",
				Context: fmt.Sprintf("%s.%s", topic, fieldName),
				Message: fmt.Sprintf(
					"@subscribe 메시지 필드 %q가 @publish 토픽 %q의 payload에 없습니다 (payload: %s)",
					fieldName, topic, joinKeys(pubFields),
				),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("@publish payload에 %q 필드를 추가하거나 @subscribe 메시지 struct에서 제거하세요", fieldName),
			})
		}
	}
	return errs
}
