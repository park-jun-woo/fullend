//ff:type feature=crosscheck type=model
//ff:what stateRef — @state 참조의 diagramID+funcName 쌍
package crosscheck

type stateRef struct {
	diagramID string
	funcName  string
}
