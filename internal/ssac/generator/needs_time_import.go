//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 파생 인터페이스에서 time 패키지 import이 필요한지 확인
package generator

func needsTimeImport(interfaces []derivedInterface) bool {
	for _, iface := range interfaces {
		if methodsHaveParamType(iface.Methods, "time.Time") {
			return true
		}
	}
	return false
}
