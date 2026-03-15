//ff:type feature=orchestrator type=model
//ff:what ChainLink represents one SSOT or artifact node in a feature chain.

package orchestrator

// ChainLink represents one SSOT or artifact node in a feature chain.
type ChainLink struct {
	Kind      string // "OpenAPI", "SSaC", "DDL", "Rego", "StateDiag", "FuncSpec", "Hurl", "STML", "Handler", "Model", "Authz", "Types"
	File      string // relative path from specs-dir or artifacts-dir
	Line      int    // 1-based line number, 0 if unknown
	Summary   string // brief description of the match
	Ownership string // "", "gen", "preserve" (empty for SSOT nodes)
}
