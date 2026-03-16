//ff:func feature=ssac-validate type=rule control=iteration dimension=2
//ff:what subscribe 함수의 message.Field 참조가 struct 필드에 존재하는지 검증한다
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateMessageFields는 subscribe 함수에서 message.Field 참조의 유효성을 검증한다.
func validateMessageFields(sf parser.ServiceFunc, seqIdx int, seq parser.Sequence) []ValidationError {
	var errs []ValidationError

	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, "message.") {
			field := val[len("message."):]
			if !hasStructField(sf.Structs, sf.Subscribe.MessageType, field) {
				ctx := errCtx{sf.FileName, sf.Name, seqIdx}
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("message.%s — 메시지 타입 %q에 %q 필드가 없습니다", field, sf.Subscribe.MessageType, field)))
			}
		}
	}

	return errs
}
