//ff:func feature=stml-gen type=generator control=sequence
//ff:what ActionBlock에 대한 useForm 훅 호출 코드를 생성한다
package generator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/stml/parser"
)

// renderFormHook generates a useForm hook call.
func renderFormHook(a parser.ActionBlock) string {
	formName := toLowerFirst(a.OperationID) + "Form"
	return fmt.Sprintf(`const %s = useForm()`, formName)
}
