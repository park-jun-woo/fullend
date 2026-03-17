//ff:func feature=crosscheck type=rule control=sequence dimension=1 topic=ssac-ddl
//ff:what 단일 시퀀스의 input key 대소문자 일치 검증
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkSeqInputKeyCase(ctx string, seqIdx int, seq ssacparser.Sequence, st *ssacvalidator.SymbolTable) []CrossError {
	if seq.Model == "" || seq.Type == "call" || seq.Package != "" {
		return nil
	}
	parts := strings.SplitN(seq.Model, ".", 2)
	if len(parts) < 2 {
		return nil
	}
	modelName, methodName := parts[0], parts[1]
	ms, ok := st.Models[modelName]
	if !ok {
		return nil
	}
	mi, exists := ms.Methods[methodName]
	if !exists || len(mi.Params) == 0 {
		return nil
	}
	return matchInputKeysToParams(ctx, seqIdx, seq, mi.Params)
}
