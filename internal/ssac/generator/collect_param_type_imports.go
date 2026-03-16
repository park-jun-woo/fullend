//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 요청 파라미터 타입에 따른 import을 수집
package generator

func collectParamTypeImports(reqParams []typedRequestParam, seen map[string]bool) {
	for _, tp := range reqParams {
		switch tp.goType {
		case "int64", "float64", "bool":
			seen["strconv"] = true
		case "time.Time":
			seen["time"] = true
		case "json.RawMessage":
			seen["encoding/json"] = true
		}
	}
}
