//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-param 바인딩이 OpenAPI parameters에 존재하는지 검증
package validator

import "github.com/geul-org/fullend/internal/stml/parser"

func validateParams(params []parser.ParamBind, opID, file string, api APISymbol) []ValidationError {
	var errs []ValidationError
	for _, p := range params {
		if !hasMatchingParam(api, p.Name) {
			errs = append(errs, errParamNotFound(file, opID, p.Name))
		}
	}
	return errs
}
