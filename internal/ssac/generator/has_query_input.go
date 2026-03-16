//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=query-opts
//ff:what Inputs map에 query 예약 소스가 있는지 확인
package generator

// hasQueryInput은 Inputs map에 query 예약 소스가 있는지 확인한다.
func hasQueryInput(inputs map[string]string) bool {
	for _, v := range inputs {
		if v == "query" {
			return true
		}
	}
	return false
}
