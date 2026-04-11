//ff:type feature=crosscheck type=model
//ff:what callRef — @call 참조의 키+컨텍스트 쌍
package crosscheck

type callRef struct {
	key     string
	context string
}
