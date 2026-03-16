//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 파라미터 배열에 특정 Go 타입이 있는지 확인
package generator

func paramsContainType(params []derivedParam, goType string) bool {
	for _, p := range params {
		if p.GoType == goType {
			return true
		}
	}
	return false
}
