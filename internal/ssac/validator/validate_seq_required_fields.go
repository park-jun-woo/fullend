//ff:func feature=ssac-validate type=rule control=selection topic=args-inputs
//ff:what 시퀀스 타입별 필수 필드 누락을 검증한다
package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateSeqRequiredFields는 단일 시퀀스의 타입별 필수 필드 누락을 검증한다.
func validateSeqRequiredFields(seq parser.Sequence, ctx errCtx) []ValidationError {
	var errs []ValidationError

	switch seq.Type {
	case parser.SeqGet:
		if seq.Model == "" {
			errs = append(errs, ctx.err("@get", "Model 누락"))
		}
		if seq.Result == nil {
			errs = append(errs, ctx.err("@get", "Result 누락"))
		}
		// Args는 0개 허용 (비즈니스 필터 없는 전체 조회)

	case parser.SeqPost:
		if seq.Model == "" {
			errs = append(errs, ctx.err("@post", "Model 누락"))
		}
		if seq.Result == nil {
			errs = append(errs, ctx.err("@post", "Result 누락"))
		}
		if len(seq.Inputs) == 0 {
			errs = append(errs, ctx.err("@post", "Inputs 누락"))
		}

	case parser.SeqPut:
		if seq.Model == "" {
			errs = append(errs, ctx.err("@put", "Model 누락"))
		}
		if seq.Result != nil {
			errs = append(errs, ctx.err("@put", "Result는 nil이어야 함"))
		}
		if len(seq.Inputs) == 0 {
			errs = append(errs, ctx.err("@put", "Inputs 누락"))
		}

	case parser.SeqDelete:
		if seq.Model == "" {
			errs = append(errs, ctx.err("@delete", "Model 누락"))
		}
		if seq.Result != nil {
			errs = append(errs, ctx.err("@delete", "Result는 nil이어야 함"))
		}
		if len(seq.Inputs) == 0 && !seq.SuppressWarn {
			errs = append(errs, ValidationError{
				FileName: ctx.fileName, FuncName: ctx.funcName, SeqIndex: ctx.seqIndex,
				Tag: "@delete", Message: "Inputs가 없습니다 — 전체 삭제 의도가 맞는지 확인하세요", Level: "WARNING",
			})
		}

	case parser.SeqEmpty, parser.SeqExists:
		if seq.Target == "" {
			errs = append(errs, ctx.err("@"+seq.Type, "Target 누락"))
		}
		if seq.Message == "" {
			errs = append(errs, ctx.err("@"+seq.Type, "Message 누락"))
		}

	case parser.SeqState:
		if seq.DiagramID == "" {
			errs = append(errs, ctx.err("@state", "DiagramID 누락"))
		}
		if len(seq.Inputs) == 0 {
			errs = append(errs, ctx.err("@state", "Inputs 누락"))
		}
		if seq.Transition == "" {
			errs = append(errs, ctx.err("@state", "Transition 누락"))
		}
		if seq.Message == "" {
			errs = append(errs, ctx.err("@state", "Message 누락"))
		}

	case parser.SeqAuth:
		if seq.Action == "" {
			errs = append(errs, ctx.err("@auth", "Action 누락"))
		}
		if seq.Resource == "" {
			errs = append(errs, ctx.err("@auth", "Resource 누락"))
		}
		if seq.Message == "" {
			errs = append(errs, ctx.err("@auth", "Message 누락"))
		}

	case parser.SeqCall:
		if seq.Model == "" {
			errs = append(errs, ctx.err("@call", "Model 누락"))
		}
		if seq.Result != nil && isPrimitiveType(seq.Result.Type) {
			errs = append(errs, ctx.err("@call", fmt.Sprintf("반환 타입에 기본 타입 %q은 사용할 수 없습니다 — Response struct 타입을 지정하세요", seq.Result.Type)))
		}

	case parser.SeqPublish:
		if seq.Topic == "" {
			errs = append(errs, ctx.err("@publish", "Topic 누락"))
		}
		if len(seq.Inputs) == 0 {
			errs = append(errs, ctx.err("@publish", "Payload 누락"))
		}

	case parser.SeqResponse:
		// Fields는 선택 — 빈 @response {} 허용 (DELETE 등)

	default:
		errs = append(errs, ctx.err("@sequence", fmt.Sprintf("알 수 없는 타입: %q", seq.Type)))
	}

	return errs
}
