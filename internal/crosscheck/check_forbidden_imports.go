//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what I/O 금지 import 목록과 대조하여 위반 패키지 반환
package crosscheck

// forbiddenImportPrefixes are DB/network packages that @call func must not import.
// File I/O (io, bufio, os) and session/cache read/write are allowed for all @call funcs.
var forbiddenImportPrefixes = []string{
	"database/sql",
	"github.com/lib/pq",
	"github.com/jackc/pgx",
	"net/http",
	"net/rpc",
	"google.golang.org/grpc",
}

// checkForbiddenImports returns any forbidden I/O imports found in the list.
func checkForbiddenImports(imports []string) []string {
	var found []string
	for _, imp := range imports {
		if isForbiddenImport(imp) {
			found = append(found, imp)
		}
	}
	return found
}
