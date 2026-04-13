//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkErrStatus — @empty/@exists/@state/@auth ErrStatus OpenAPI 응답 존재 검증 (X-21, X-22)
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/fullend"

func checkErrStatus(fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if fn.Subscribe != nil {
			continue
		}
		errs = append(errs, checkFuncErrStatus(fn.Name, fn.Sequences, fs)...)
	}
	return errs
}
