//ff:func feature=ssac-validate type=rule
//ff:what subscribe/HTTP 트리거 관련 규칙 검증

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateSubscribeRules는 subscribe/HTTP 트리거와 관련된 규칙을 검증한다.
func validateSubscribeRules(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError

	if sf.Subscribe != nil {
		ctx := errCtx{sf.FileName, sf.Name, -1}

		// subscribe 함수에 파라미터 필수
		if sf.Param == nil {
			errs = append(errs, ctx.err("@subscribe", "@subscribe 함수에 파라미터가 필요합니다 — func Name(TypeName message) {}"))
		}

		// 파라미터 변수명은 반드시 "message"
		if sf.Param != nil && sf.Param.VarName != "message" {
			errs = append(errs, ctx.err("@subscribe", fmt.Sprintf("파라미터 변수명은 \"message\"여야 합니다 — 현재: %q", sf.Param.VarName)))
		}

		// MessageType이 파일 내 struct로 존재하는지
		if sf.Subscribe.MessageType != "" {
			found := false
			for _, si := range sf.Structs {
				if si.Name == sf.Subscribe.MessageType {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, ctx.err("@subscribe", fmt.Sprintf("메시지 타입 %q이 파일 내에 struct로 선언되지 않았습니다", sf.Subscribe.MessageType)))
			}
		}

		// subscribe 함수에 @response 있으면 ERROR
		for i, seq := range sf.Sequences {
			if seq.Type == parser.SeqResponse {
				seqCtx := errCtx{sf.FileName, sf.Name, i}
				errs = append(errs, seqCtx.err("@subscribe", "@subscribe 함수에 @response를 사용할 수 없습니다"))
			}
		}
		// subscribe 함수에서 request 사용하면 ERROR
		for i, seq := range sf.Sequences {
			for _, val := range seq.Inputs {
				if strings.HasPrefix(val, "request.") {
					seqCtx := errCtx{sf.FileName, sf.Name, i}
					errs = append(errs, seqCtx.err("@subscribe", "@subscribe 함수에서 request를 사용할 수 없습니다 — message를 사용하세요"))
					break
				}
			}
		}
		// subscribe 함수에서 query 사용하면 ERROR
		for i, seq := range sf.Sequences {
			for _, val := range seq.Inputs {
				if val == "query" || strings.HasPrefix(val, "query.") {
					seqCtx := errCtx{sf.FileName, sf.Name, i}
					errs = append(errs, seqCtx.err("@subscribe", "query는 HTTP 전용입니다 — @subscribe 함수에서 사용할 수 없습니다"))
					break
				}
			}
		}

		// message.Field 검증: struct 필드 존재 확인
		if sf.Subscribe.MessageType != "" {
			for i, seq := range sf.Sequences {
				for _, val := range seq.Inputs {
					if strings.HasPrefix(val, "message.") {
						field := val[len("message."):]
						if !hasStructField(sf.Structs, sf.Subscribe.MessageType, field) {
							seqCtx := errCtx{sf.FileName, sf.Name, i}
							errs = append(errs, seqCtx.err("@"+seq.Type, fmt.Sprintf("message.%s — 메시지 타입 %q에 %q 필드가 없습니다", field, sf.Subscribe.MessageType, field)))
						}
					}
				}
			}
		}
	} else {
		// HTTP 함수에서 message 사용하면 ERROR
		for i, seq := range sf.Sequences {
			for _, val := range seq.Inputs {
				if strings.HasPrefix(val, "message.") {
					ctx := errCtx{sf.FileName, sf.Name, i}
					errs = append(errs, ctx.err("@sequence", "HTTP 함수에서 message를 사용할 수 없습니다 — @subscribe 함수에서만 사용 가능합니다"))
					break
				}
			}
		}
	}

	return errs
}
