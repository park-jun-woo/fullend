//ff:type feature=orchestrator type=model
//ff:what StatusLine holds one SSOT's status info.

package orchestrator

// StatusLine holds one SSOT's status info.
type StatusLine struct {
	Kind    SSOTKind
	Path    string // relative path for display
	Summary string
}
