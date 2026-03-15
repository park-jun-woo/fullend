//ff:type feature=symbol type=model
//ff:what API 엔드포인트의 request/response 필드 목록 + HasQueryOpts
package validator

// OperationSymbol은 API 엔드포인트의 request/response 필드 목록이다.
type OperationSymbol struct {
	RequestFields map[string]bool
	PathParams    []PathParam // path parameter (순서 보존)
	XPagination    *XPagination
	XSort          *XSort
	XFilter        *XFilter
	XInclude       *XInclude
}

// HasQueryOpts는 x- 확장이 하나라도 있는지 반환한다.
func (op OperationSymbol) HasQueryOpts() bool {
	return op.XPagination != nil || op.XSort != nil || op.XFilter != nil || op.XInclude != nil
}
