//ff:func feature=gen-gogin type=util
//ff:what generates the struct field assignment for VerifyToken result

package gogin

import (
	"fmt"

	"github.com/geul-org/fullend/internal/projectconfig"
)

// resultAssignLine generates the struct field assignment for VerifyToken result.
func resultAssignLine(field string, def projectconfig.ClaimDef) string {
	varName := lcFirst(field) + "Raw"
	switch def.GoType {
	case "int64":
		return fmt.Sprintf("\t\t%s: int64(%s),", field, varName)
	default:
		return fmt.Sprintf("\t\t%s: %s,", field, varName)
	}
}
