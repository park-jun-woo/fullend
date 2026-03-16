//ff:type feature=genapi type=model
//ff:what 백엔드 코드 생성 인터페이스
package genapi

// Backend generates backend code from parsed SSOTs.
type Backend interface {
	Generate(parsed *ParsedSSOTs, cfg *GenConfig) error
}
