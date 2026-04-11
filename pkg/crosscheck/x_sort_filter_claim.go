//ff:type feature=crosscheck type=model
//ff:what xSortFilterClaim — x-sort/x-filter 검증 대상
package crosscheck

type xSortFilterClaim struct {
	ruleID    string
	col       string
	lookupKey string
	context   string
	message   string
}
