//ff:type feature=ssac-gen type=model topic=request-params
//ff:what 미가공 요청 파라미터(이름+Go타입)를 담는 구조체
package generator

type rawParam struct {
	name   string
	goType string
}
