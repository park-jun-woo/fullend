//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-paginate·sort·filter 인프라 파라미터 검증
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func validateInfraParams(f parser.FetchBlock, file string, api APISymbol) []ValidationError {
	var errs []ValidationError

	if f.Paginate && api.Pagination == nil {
		errs = append(errs, errPaginateNoExt(file, f.OperationID))
	}

	if f.Sort != nil {
		if api.Sort == nil {
			errs = append(errs, errSortNotAllowed(file, f.OperationID, f.Sort.Column))
		} else if !containsStr(api.Sort.Allowed, f.Sort.Column) {
			errs = append(errs, errSortNotAllowed(file, f.OperationID, f.Sort.Column))
		}
	}

	for _, col := range f.Filters {
		if api.Filter == nil || !containsStr(api.Filter.Allowed, col) {
			errs = append(errs, errFilterNotAllowed(file, f.OperationID, col))
		}
	}

	return errs
}
