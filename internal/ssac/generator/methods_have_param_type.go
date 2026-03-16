//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 메서드 배열의 파라미터에 특정 Go 타입이 있는지 확인
package generator

func methodsHaveParamType(methods []derivedMethod, goType string) bool {
	for _, m := range methods {
		if paramsContainType(m.Params, goType) {
			return true
		}
	}
	return false
}
