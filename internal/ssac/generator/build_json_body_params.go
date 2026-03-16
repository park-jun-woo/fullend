//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what 여러 요청 파라미터를 JSON body 바인딩 코드로 변환
package generator

import (
	"bytes"
	"fmt"

	"github.com/ettle/strcase"
)

func buildJSONBodyParams(rawParams []rawParam) []typedRequestParam {
	var buf bytes.Buffer

	buf.WriteString("\tvar req struct {\n")
	for _, rp := range rawParams {
		buf.WriteString(fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", strcase.ToGoPascal(rp.name), rp.goType, rp.name))
	}
	buf.WriteString("\t}\n")
	buf.WriteString("\tif err := c.ShouldBindJSON(&req); err != nil {\n")
	buf.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"invalid request body\"})\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	for _, rp := range rawParams {
		varName := strcase.ToGoCamel(rp.name)
		buf.WriteString(fmt.Sprintf("\t%s := req.%s\n", varName, strcase.ToGoPascal(rp.name)))
	}

	result := []typedRequestParam{{
		name:        "_json_body",
		goType:      "json_body",
		extractCode: buf.String(),
	}}
	for _, rp := range rawParams {
		if rp.goType == "time.Time" || rp.goType == "json.RawMessage" {
			result = append(result, typedRequestParam{name: rp.name, goType: rp.goType})
		}
	}
	return result
}
