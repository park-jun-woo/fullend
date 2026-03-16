//ff:func feature=gen-gogin type=util control=selection
//ff:what converts a numeric HTTP status code string to Go net/http constant name

package gogin

// httpStatusConst converts a numeric HTTP status code string to Go's net/http constant name.
func httpStatusConst(code string) string {
	switch code {
	case "200":
		return "http.StatusOK"
	case "201":
		return "http.StatusCreated"
	case "204":
		return "http.StatusNoContent"
	default:
		return "http.StatusOK"
	}
}
