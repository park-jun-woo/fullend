//ff:func feature=genmodel type=generator control=iteration dimension=1
//ff:what 경로 파라미터만 필터링하여 반환한다
package genmodel

func (m methodInfo) pathParams() []paramInfo {
	var result []paramInfo
	for _, p := range m.Params {
		if p.In == "path" {
			result = append(result, p)
		}
	}
	return result
}
