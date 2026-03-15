//ff:func feature=orchestrator type=util
//ff:what KindFromString parses a CLI --skip value into a SSOTKind.

package orchestrator

import "strings"

// kindNames maps CLI --skip values to SSOTKind.
var kindNames = map[string]SSOTKind{
	"openapi":  KindOpenAPI,
	"ddl":      KindDDL,
	"ssac":     KindSSaC,
	"model":    KindModel,
	"stml":     KindSTML,
	"states":   KindStates,
	"policy":   KindPolicy,
	"scenario": KindScenario,
	"func":     KindFunc,
	"config":   KindConfig,
}

// KindFromString parses a CLI --skip value into a SSOTKind.
func KindFromString(s string) (SSOTKind, bool) {
	k, ok := kindNames[strings.ToLower(s)]
	return k, ok
}
