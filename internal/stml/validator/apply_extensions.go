//ff:func feature=stml-validate type=parser control=sequence
//ff:what x-pagination·x-sort·x-filter 확장을 APISymbol에 적용
package validator

// applyExtensions applies x-pagination, x-sort, x-filter extensions to an APISymbol.
func applyExtensions(op openAPIOperation, api *APISymbol) {
	if op.XPagination != nil {
		api.Pagination = &PaginationExt{
			Style:        op.XPagination.Style,
			DefaultLimit: op.XPagination.DefaultLimit,
			MaxLimit:     op.XPagination.MaxLimit,
		}
	}
	if op.XSort != nil {
		dir := op.XSort.Direction
		if dir == "" {
			dir = "asc"
		}
		api.Sort = &SortExt{
			Allowed:   op.XSort.Allowed,
			Default:   op.XSort.Default,
			Direction: dir,
		}
	}
	if op.XFilter != nil {
		api.Filter = &FilterExt{Allowed: op.XFilter.Allowed}
	}
}
