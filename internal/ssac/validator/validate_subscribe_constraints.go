//ff:func feature=ssac-validate type=rule control=iteration dimension=2
//ff:what subscribe 함수의 파라미터·타입·시퀀스 제약 조건을 검증한다
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateSubscribeConstraints는 subscribe 함수의 제약 조건을 검증한다.
func validateSubscribeConstraints(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
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
			errs = append(errs, validateMessageFields(sf, i, seq)...)
		}
	}

	return errs
}
