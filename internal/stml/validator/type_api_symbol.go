//ff:type feature=stml-validate type=model
//ff:what 단일 OpenAPI 오퍼레이션을 나타내는 심볼
package validator

// APISymbol represents a single OpenAPI operation.
type APISymbol struct {
	Method         string                 // "get", "post", "put", "delete"
	Parameters     []ParamSymbol          // path/query parameters
	RequestFields  map[string]string      // field name → type
	ResponseFields map[string]FieldSymbol // field name → type info

	// Phase 5: x- extensions
	Pagination *PaginationExt
	Sort       *SortExt
	Filter     *FilterExt
}
