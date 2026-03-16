//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what operationId에서 CRUD 접두사(List/Get/Create/Update/Delete) 제거
package crosscheck

import "strings"

// stripCRUDPrefix removes CRUD prefixes from an operation name.
func stripCRUDPrefix(name string) string {
	for _, prefix := range []string{"List", "Get", "Create", "Update", "Delete"} {
		name = strings.TrimPrefix(name, prefix)
	}
	return name
}
