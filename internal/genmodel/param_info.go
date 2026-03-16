//ff:type feature=genmodel type=model
//ff:what 파라미터 정보 타입 정의
package genmodel

type paramInfo struct {
	Name   string
	GoType string
	In     string // "body", "path"
}
