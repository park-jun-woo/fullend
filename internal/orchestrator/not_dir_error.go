//ff:type feature=orchestrator type=model
//ff:what NotDirError is returned when the specs path is not a directory.

package orchestrator

// NotDirError is returned when the specs path is not a directory.
type NotDirError struct {
	Path string
}
