//ff:type feature=reporter type=model
//ff:what 기능 체인의 링크를 나타내는 구조체
package reporter

// ChainLink mirrors orchestrator.ChainLink to avoid circular import.
type ChainLink struct {
	Kind      string
	File      string
	Line      int
	Summary   string
	Ownership string // "", "gen", "preserve"
}
