//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkUndeclaredSources — @call arg source가 선언되지 않은 변수인지 WARNING
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func checkUndeclaredSources(funcName string, seqs []ssac.Sequence, declared map[string]bool) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		errs = append(errs, checkUndeclaredArgs(funcName, seq.Args, declared)...)
	}
	return errs
}
