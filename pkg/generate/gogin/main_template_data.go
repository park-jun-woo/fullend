//ff:type feature=gen-gogin type=model topic=output
//ff:what MainTmplData — main.tmpl 렌더링 입력 struct

package gogin

// MainTmplData carries named fields for rendering templates/main.tmpl.
// 모든 필드는 사전 조립된 문자열 (조건부 블록은 caller 가 빈 문자열로 넘김).
type MainTmplData struct {
	OsImport            string
	ImportBlock         string
	QueueImport         string
	BuiltinImport       string
	JWTFlagLine         string
	AuthzBlock          string
	QueueInitBlock      string
	BuiltinInitBlock    string
	InitBlock           string
	QueueSubscribeBlock string
	DBName              string
}
