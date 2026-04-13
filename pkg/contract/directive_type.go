//ff:type feature=contract type=model
//ff:what fullend 소유권 디렉티브를 나타내는 구조체
package contract

// Directive represents a //fullend: ownership directive attached to generated code.
type Directive struct {
	Ownership string // "gen" or "preserve"
	SSOT      string // SSOT file relative path (e.g. "service/gig/create_gig.ssac")
	Contract  string // 7-char SHA256 hex prefix
}
