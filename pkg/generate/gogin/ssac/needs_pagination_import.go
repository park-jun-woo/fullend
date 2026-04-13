//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 파생 인터페이스에서 pagination 패키지 import이 필요한지 확인
package ssac

func needsPaginationImport(interfaces []derivedInterface) bool {
	for _, iface := range interfaces {
		if methodsHaveReturnSubstring(iface.Methods, "pagination.") {
			return true
		}
	}
	return false
}
