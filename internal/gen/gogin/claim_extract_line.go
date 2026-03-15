//ff:func feature=gen-gogin type=util
//ff:what generates the JWT MapClaims extraction line for VerifyToken

package gogin

import (
	"fmt"

	"github.com/geul-org/fullend/internal/projectconfig"
)

// claimExtractLine generates the JWT MapClaims extraction line for VerifyToken.
func claimExtractLine(field string, def projectconfig.ClaimDef) string {
	varName := lcFirst(field) + "Raw"
	switch def.GoType {
	case "int64":
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(float64)", varName, def.Key)
	case "bool":
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(bool)", varName, def.Key)
	default: // string
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(string)", varName, def.Key)
	}
}
