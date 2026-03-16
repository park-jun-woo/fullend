//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 오퍼레이션의 parameters를 APISymbol에 수집
package validator

func collectParams(op openAPIOperation, api *APISymbol) {
	for _, p := range op.Parameters {
		api.Parameters = append(api.Parameters, ParamSymbol{
			Name: p.Name,
			In:   p.In,
		})
	}
}
