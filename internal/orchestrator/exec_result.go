//ff:type feature=orchestrator type=model
//ff:what ExecResult holds the outcome of an external command execution.

package orchestrator

// ExecResult holds the outcome of an external command execution.
type ExecResult struct {
	Skipped bool   // true if the tool is not installed
	Err     error  // non-nil if the command failed
	Stderr  string // captured stderr on failure
}
