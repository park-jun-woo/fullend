//ff:func feature=ssac-gen type=generator control=selection topic=request-params
//ff:what Go 타입별 요청 파라미터 추출 코드를 생성
package ssac

import "fmt"

func generateExtractCode(varName, paramName, goType string) string {
	switch goType {
	case "int64":
		return fmt.Sprintf("\t%s, err := strconv.ParseInt(c.Query(%q), 10, 64)\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s: 유효하지 않은 값\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "float64":
		return fmt.Sprintf("\t%s, err := strconv.ParseFloat(c.Query(%q), 64)\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s: 유효하지 않은 값\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "bool":
		return fmt.Sprintf("\t%s, err := strconv.ParseBool(c.Query(%q))\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s: 유효하지 않은 값\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	case "time.Time":
		return fmt.Sprintf("\t%s, err := time.Parse(time.RFC3339, c.Query(%q))\n"+
			"\tif err != nil {\n"+
			"\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s: 유효하지 않은 값\"})\n"+
			"\t\treturn\n"+
			"\t}\n", varName, paramName, paramName)
	default:
		return fmt.Sprintf("\t%s := c.Query(%q)\n", varName, paramName)
	}
}
