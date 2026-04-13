//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=request-params
//ff:what 여러 요청 파라미터를 JSON body 바인딩 코드로 변환
package ssac

import (
	"bytes"
	"fmt"

	"github.com/ettle/strcase"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func buildJSONBodyParams(rawParams []rawParam, rs *rule.RequestSchemaInfo) []typedRequestParam {
	var buf bytes.Buffer

	buf.WriteString("\tvar req struct {\n")
	for _, rp := range rawParams {
		tag := buildFieldTag(rp.name, rs)
		buf.WriteString(fmt.Sprintf("\t\t%s %s `%s`\n", strcase.ToGoPascal(rp.name), rp.goType, tag))
	}
	buf.WriteString("\t}\n")
	buf.WriteString("\tif err := c.ShouldBindJSON(&req); err != nil {\n")
	buf.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})\n")
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
