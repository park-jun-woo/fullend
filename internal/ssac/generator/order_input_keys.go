//ff:func feature=ssac-gen type=util control=selection
//ff:what paramOrder 유무에 따라 입력 키 정렬 전략을 선택
package generator

func orderInputKeys(inputs map[string]string, paramOrder []string) []string {
	switch {
	case len(paramOrder) > 0:
		return orderByParamOrder(inputs, paramOrder)
	default:
		return orderAlphabetQueryLast(inputs)
	}
}
