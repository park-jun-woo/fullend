//ff:type feature=orchestrator type=model
//ff:what SSOTKind identifies a type of SSOT source with its const variants.

package orchestrator

// SSOTKind identifies a type of SSOT source.
type SSOTKind string

const (
	KindOpenAPI  SSOTKind = "OpenAPI"
	KindDDL      SSOTKind = "DDL"
	KindSSaC     SSOTKind = "SSaC"
	KindModel    SSOTKind = "Model"
	KindSTML     SSOTKind = "STML"
	KindStates   SSOTKind = "States"
	KindPolicy   SSOTKind = "Policy"
	KindScenario SSOTKind = "Scenario"
	KindFunc     SSOTKind = "Func"
	KindConfig   SSOTKind = "Config"
)
