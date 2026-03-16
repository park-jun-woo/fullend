//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 미구현 func의 스켈레톤 코드 힌트 생성
package crosscheck

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// generateSkeleton creates a skeleton code hint for a missing func.
func generateSkeleton(pkg, funcName string, seq ssacparser.Sequence) string {
	uc := strcase.ToGoPascal(funcName)
	if pkg == "" {
		pkg = "custom"
	}

	var requestFields []string
	for key := range seq.Inputs {
		requestFields = append(requestFields, fmt.Sprintf("\t%s string", key))
	}

	var responseFields []string
	if seq.Result != nil {
		typeName := "string"
		if seq.Result.Type != "" {
			typeName = seq.Result.Type
		}
		responseFields = append(responseFields, fmt.Sprintf("\t%s %s", strcase.ToGoPascal(seq.Result.Var), typeName))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("다음 파일을 작성하세요: func/%s/%s.go\n\n", pkg, toSnakeCase(funcName)))
	b.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	b.WriteString(fmt.Sprintf("// @func %s\n", funcName))
	b.WriteString("// @description <이 함수가 무엇을 하는지 한 줄로 설명>\n\n")
	b.WriteString(fmt.Sprintf("type %sRequest struct {\n", uc))
	for _, f := range requestFields {
		b.WriteString(f + "\n")
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("type %sResponse struct {\n", uc))
	for _, f := range responseFields {
		b.WriteString(f + "\n")
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("func %s(req %sRequest) (%sResponse, error) {\n", uc, uc, uc))
	b.WriteString("\t// TODO: implement\n")
	b.WriteString(fmt.Sprintf("\treturn %sResponse{}, nil\n", uc))
	b.WriteString("}\n")

	return b.String()
}
