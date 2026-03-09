package crosscheck

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/artifacts/internal/funcspec"
	ssacparser "github.com/geul-org/ssac/parser"
)

// CheckFuncs validates SSaC @func references against parsed func specs.
// fullendPkgSpecs: specs from pkg/ (fullend default).
// projectFuncSpecs: specs from specs/<project>/func/ (custom).
func CheckFuncs(serviceFuncs []ssacparser.ServiceFunc, fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec) []CrossError {
	var errs []CrossError

	// Build lookup: "package.funcName" → FuncSpec.
	// Project custom overrides fullend default.
	specMap := make(map[string]*funcspec.FuncSpec)
	for i := range fullendPkgSpecs {
		key := fullendPkgSpecs[i].Package + "." + fullendPkgSpecs[i].Name
		specMap[key] = &fullendPkgSpecs[i]
	}
	for i := range projectFuncSpecs {
		key := projectFuncSpecs[i].Package + "." + projectFuncSpecs[i].Name
		specMap[key] = &projectFuncSpecs[i]
	}

	// Collect SSaC @func references.
	for _, sf := range serviceFuncs {
		for seqIdx, seq := range sf.Sequences {
			if seq.Type != "call" || seq.Func == "" {
				continue
			}

			pkg := seq.Package
			funcName := seq.Func
			key := pkg + "." + funcName
			if pkg == "" {
				key = funcName
			}

			spec, found := specMap[key]
			if !found {
				// Generate skeleton hint.
				skeleton := generateSkeleton(pkg, funcName, seq)
				errs = append(errs, CrossError{
					Rule:       "Func ↔ SSaC",
					Context:    fmt.Sprintf("%s seq[%d] @func %s", sf.Name, seqIdx, key),
					Message:    fmt.Sprintf("@func %s — 구현 없음", key),
					Level:      "ERROR",
					Suggestion: skeleton,
				})
				continue
			}

			// Check HasBody.
			if !spec.HasBody {
				errs = append(errs, CrossError{
					Rule:    "Func ↔ SSaC",
					Context: fmt.Sprintf("%s seq[%d] @func %s", sf.Name, seqIdx, key),
					Message: "본체 미구현 (TODO)",
					Level:   "WARNING",
				})
			}
		}
	}

	return errs
}

// generateSkeleton creates a skeleton code hint for a missing func.
func generateSkeleton(pkg, funcName string, seq ssacparser.Sequence) string {
	uc := strings.ToUpper(funcName[:1]) + funcName[1:]
	if pkg == "" {
		pkg = "custom"
	}

	var inputFields []string
	for _, p := range seq.Params {
		name := p.Name
		if strings.HasPrefix(name, "\"") {
			continue // literal
		}
		if p.Source == "request" {
			inputFields = append(inputFields, fmt.Sprintf("\t%s string", name))
		} else if strings.Contains(name, ".") {
			parts := strings.SplitN(name, ".", 2)
			inputFields = append(inputFields, fmt.Sprintf("\t%s string", parts[1]))
		}
	}

	var outputFields []string
	if seq.Result != nil {
		typeName := "string"
		if seq.Result.Type != "" {
			typeName = seq.Result.Type
		}
		outputFields = append(outputFields, fmt.Sprintf("\t%s %s", strings.ToUpper(seq.Result.Var[:1])+seq.Result.Var[1:], typeName))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("다음 파일을 작성하세요: func/%s/%s.go\n\n", pkg, toSnakeCase(funcName)))
	b.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	b.WriteString(fmt.Sprintf("// @func %s\n", funcName))
	b.WriteString("// @description <이 함수가 무엇을 하는지 한 줄로 설명>\n\n")
	b.WriteString(fmt.Sprintf("type %sInput struct {\n", uc))
	for _, f := range inputFields {
		b.WriteString(f + "\n")
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("type %sOutput struct {\n", uc))
	for _, f := range outputFields {
		b.WriteString(f + "\n")
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("func %s(in %sInput) (%sOutput, error) {\n", uc, uc, uc))
	b.WriteString("\t// TODO: implement\n")
	b.WriteString(fmt.Sprintf("\treturn %sOutput{}, nil\n", uc))
	b.WriteString("}\n")

	return b.String()
}

// toSnakeCase converts camelCase to snake_case.
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c-'A'+'a'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
