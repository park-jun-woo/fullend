//ff:func feature=ssac-gen type=generator control=selection topic=path-params
//ff:what 경로 파라미터의 Go 추출 코드를 생성
package ssac

import (
	"fmt"

	"github.com/ettle/strcase"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func generatePathParamCode(pp validator.PathParam) string {
	varName := strcase.ToGoCamel(pp.Name)
	switch pp.GoType {
	case "int64":
		return fmt.Sprintf("\t%sStr := c.Param(%q)\n"+
			"\t%s, err := strconv.ParseInt(%sStr, 10, 64)\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"invalid path parameter\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, pp.Name, varName, varName)
	case "float64":
		return fmt.Sprintf("\t%sStr := c.Param(%q)\n"+
			"\t%s, err := strconv.ParseFloat(%sStr, 64)\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"invalid path parameter\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, pp.Name, varName, varName)
	default:
		return fmt.Sprintf("\t%s := c.Param(%q)\n", varName, pp.Name)
	}
}
