//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 파생 인터페이스에서 encoding/json 패키지 import이 필요한지 확인
package ssac

func needsJSONImport(interfaces []derivedInterface) bool {
	for _, iface := range interfaces {
		if methodsHaveParamType(iface.Methods, "json.RawMessage") {
			return true
		}
	}
	return false
}
