//ff:func feature=gen-gogin type=util control=selection
//ff:what generates the JWT MapClaims extraction line for VerifyToken

package gogin

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

// claimExtractLine generates the JWT MapClaims extraction line for VerifyToken.
func claimExtractLine(field string, def manifest.ClaimDef) string {
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
