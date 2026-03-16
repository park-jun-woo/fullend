//ff:func feature=ssac-gen type=util control=sequence topic=interface-derive
//ff:what 메서드의 파라미터를 시그니처 문자열로 결합 (QueryOpts 포함)
package generator

func renderMethodSignature(m derivedMethod) string {
	params := renderParams(m.Params)
	if !m.HasQueryOpts {
		return params
	}
	if params != "" {
		params += ", "
	}
	return params + "opts QueryOpts"
}
