//ff:type feature=symbol type=model
//ff:what 모델의 메서드 목록 + HasMethod
package validator

// ModelSymbol은 모델의 메서드 목록이다.
type ModelSymbol struct {
	Methods map[string]MethodInfo
}

// HasMethod는 메서드 존재 여부를 반환한다.
func (ms ModelSymbol) HasMethod(name string) bool {
	_, ok := ms.Methods[name]
	return ok
}
