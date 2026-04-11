//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSingleFuncPurity — 단일 func의 TODO/금지 import 검사
package crosscheck

var forbiddenImports = map[string]bool{
	"database/sql":           true,
	"github.com/lib/pq":     true,
	"github.com/jackc/pgx":  true,
	"net/http":               true,
	"net/rpc":                true,
	"google.golang.org/grpc": true,
}

func checkSingleFuncPurity(pkg, name string, hasBody bool, imports []string) []CrossError {
	var errs []CrossError
	key := pkg + "." + name
	if !hasBody {
		errs = append(errs, CrossError{Rule: "X-40", Context: key, Level: "ERROR", Message: "func body is TODO stub"})
	}
	for _, imp := range imports {
		if forbiddenImports[imp] {
			errs = append(errs, CrossError{Rule: "X-41", Context: key, Level: "ERROR", Message: "func imports forbidden package: " + imp})
		}
	}
	return errs
}
