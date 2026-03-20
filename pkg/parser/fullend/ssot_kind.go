//ff:type feature=orchestrator type=model
//ff:what SSOT 종류를 나타내는 열거 타입
package fullend

// SSOTKind identifies one of the SSOT types in a fullend project.
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
	KindToulmin  SSOTKind = "Toulmin"
)

// DetectedSSOT represents a detected SSOT file or directory.
type DetectedSSOT struct {
	Kind SSOTKind
	Path string // absolute path to the relevant directory or file
}
