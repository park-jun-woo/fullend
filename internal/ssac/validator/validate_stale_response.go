//ff:func feature=ssac-validate type=rule control=iteration dimension=1
//ff:what put/delete 이후 갱신 없이 response에서 사용되는 변수 경고

package validator

import "github.com/geul-org/fullend/internal/ssac/parser"

// validateStaleResponse는 put/delete 이후 갱신 없이 response에서 사용되는 변수를 경고한다.
func validateStaleResponse(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError

	getVars := map[string]string{}   // var → model
	mutated := map[string]bool{}     // model → mutated?

	for i, seq := range sf.Sequences {
		if seq.Type == parser.SeqGet {
			trackGetVar(seq, getVars, mutated)
			continue
		}
		if seq.Type == parser.SeqPut || seq.Type == parser.SeqDelete {
			trackMutation(seq, mutated)
			continue
		}
		if seq.Type != parser.SeqResponse || seq.SuppressWarn {
			continue
		}
		ctx := errCtx{sf.FileName, sf.Name, i}
		errs = collectStaleErrors(seq, getVars, mutated, ctx, errs)
	}

	return errs
}
