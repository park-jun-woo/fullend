//ff:func feature=gen-hurl type=util control=selection
//ff:what Returns sort order for CRUD methods (POST, GET, PUT, DELETE).
package hurl

func crudOrder(method string) int {
	switch method {
	case "POST":
		return 0
	case "GET":
		return 1
	case "PUT":
		return 2
	case "DELETE":
		return 3
	default:
		return 4
	}
}
