//ff:func feature=genmodel type=util control=iteration dimension=1
//ff:what 구조체 필드 중 time.Time 타입이 있는지 확인한다
package genmodel

func structHasTimeField(t structType) bool {
	for _, f := range t.Fields {
		if f.GoType == "time.Time" {
			return true
		}
	}
	return false
}
