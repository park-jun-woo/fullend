//ff:type feature=orchestrator type=model
//ff:what DetectedSSOT holds the kind and resolved directory path.

package orchestrator

// DetectedSSOT holds the kind and resolved directory path.
type DetectedSSOT struct {
	Kind SSOTKind
	Path string // absolute path to the relevant directory or file
}
