//ff:func feature=orchestrator type=model control=sequence
//ff:what NotDirError.Error returns the error message for NotDirError.

package orchestrator

func (e *NotDirError) Error() string {
	return "not a directory: " + e.Path
}
