package crosscheck

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckSSaCOpenAPI validates SSaC function names match OpenAPI operationIds and vice versa,
// and SSaC @response fields match OpenAPI response schema properties.
func CheckSSaCOpenAPI(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, doc *openapi3.T, funcSpecs []funcspec.FuncSpec) []CrossError {
	var errs []CrossError

	funcNames := make(map[string]string) // funcName → fileName
	for _, fn := range funcs {
		funcNames[fn.Name] = fn.FileName
	}

	// Rule 3: Every SSaC function must have a matching operationId.
	for name, fileName := range funcNames {
		if _, ok := st.Operations[name]; !ok {
			errs = append(errs, CrossError{
				Rule:       "SSaC → OpenAPI",
				Context:    fmt.Sprintf("%s:%s", fileName, name),
				Message:    fmt.Sprintf("SSaC function %q has no matching OpenAPI operationId", name),
				Suggestion: fmt.Sprintf("OpenAPI에 추가: operationId: %s", name),
			})
		}
	}

	// Rule 4: Every operationId should have a matching SSaC function.
	for opID := range st.Operations {
		if _, ok := funcNames[opID]; !ok {
			errs = append(errs, CrossError{
				Rule:       "OpenAPI → SSaC",
				Context:    opID,
				Message:    fmt.Sprintf("OpenAPI operationId %q has no matching SSaC function", opID),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("SSaC에 추가: func %s(w http.ResponseWriter, r *http.Request) {}", opID),
			})
		}
	}

	// Rule 5 & 6: SSaC @response fields ↔ OpenAPI response schema properties.
	errs = append(errs, checkResponseFields(funcs, st, doc, funcSpecs)...)

	// Rule 7: SSaC ErrStatus → OpenAPI error response defined.
	if doc != nil {
		errs = append(errs, checkErrStatus(funcs, doc)...)
	}

	// Rule 8: SSaC @response → OpenAPI must have explicit 2xx response code.
	if doc != nil {
		errs = append(errs, checkResponseSuccessCode(funcs, doc)...)
	}

	return errs
}

// checkResponseFields validates that SSaC @response field keys match OpenAPI response schema properties.
func checkResponseFields(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, doc *openapi3.T, funcSpecs []funcspec.FuncSpec) []CrossError {
	var errs []CrossError

	// Build OpenAPI response properties per operationId.
	opResponseProps := buildOperationResponseProps(doc)

	for _, fn := range funcs {
		// Find @response sequence with explicit fields.
		responseFields := extractResponseFieldKeys(fn)

		if responseFields == nil {
			// shorthand (@response varName) — resolve fields from variable type.
			shorthandFields := resolveShorthandResponseFields(fn, funcSpecs, st)
			if shorthandFields == nil {
				continue // no @response or type tracking failed — skip
			}

			opProps, hasOp := opResponseProps[fn.Name]
			if !hasOp {
				continue
			}

			shorthandSet := make(map[string]bool, len(shorthandFields))
			for _, f := range shorthandFields {
				shorthandSet[f] = true
			}

			// Rule 5: shorthand field → OpenAPI property (ERROR).
			for _, jf := range shorthandFields {
				if !opProps[jf] {
					errs = append(errs, CrossError{
						Rule:       "SSaC @response → OpenAPI",
						Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
						Message:    fmt.Sprintf("shorthand @response 변수의 JSON 필드 %q가 OpenAPI %s 응답 스키마에 없습니다", jf, fn.Name),
						Suggestion: fmt.Sprintf("OpenAPI %s 응답 스키마의 property명을 %q로 변경하세요", fn.Name, jf),
					})
				}
			}

			// Rule 6: OpenAPI property → shorthand field (WARNING).
			for prop := range opProps {
				if !shorthandSet[prop] {
					errs = append(errs, CrossError{
						Rule:       "OpenAPI → SSaC @response",
						Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
						Message:    fmt.Sprintf("OpenAPI %s 응답 필드 %q가 shorthand @response 변수 타입에 없습니다", fn.Name, prop),
						Level:      "WARNING",
						Suggestion: fmt.Sprintf("OpenAPI에서 %q를 제거하거나 변수 타입에 해당 필드를 추가하세요", prop),
					})
				}
			}
			continue
		}

		opProps, hasOp := opResponseProps[fn.Name]
		if !hasOp {
			continue // no OpenAPI operation — already caught by Rule 3
		}

		// Rule 5: SSaC @response field → OpenAPI response property (ERROR).
		for _, field := range responseFields {
			if !opProps[field] {
				errs = append(errs, CrossError{
					Rule:       "SSaC @response → OpenAPI",
					Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
					Message:    fmt.Sprintf("SSaC @response 필드 %q가 OpenAPI %s 응답 스키마에 없습니다", field, fn.Name),
					Suggestion: fmt.Sprintf("OpenAPI %s 응답 스키마에 %q property를 추가하세요", fn.Name, field),
				})
			}
		}

		// Rule 6: OpenAPI response property → SSaC @response field (WARNING).
		responseFieldSet := make(map[string]bool, len(responseFields))
		for _, f := range responseFields {
			responseFieldSet[f] = true
		}
		for prop := range opProps {
			if !responseFieldSet[prop] {
				errs = append(errs, CrossError{
					Rule:       "OpenAPI → SSaC @response",
					Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
					Message:    fmt.Sprintf("OpenAPI %s 응답 필드 %q가 SSaC @response에 없습니다", fn.Name, prop),
					Level:      "WARNING",
					Suggestion: fmt.Sprintf("SSaC @response에 %q 필드를 추가하거나 OpenAPI에서 제거하세요", prop),
				})
			}
		}
	}

	return errs
}

// resolveShorthandResponseFields resolves the JSON field names for a shorthand @response varName.
// It traces the variable back to its origin (@call or @get/@put/@post/@delete) and returns
// the JSON field names that would be serialized.
func resolveShorthandResponseFields(
	fn ssacparser.ServiceFunc,
	funcSpecs []funcspec.FuncSpec,
	st *ssacvalidator.SymbolTable,
) []string {
	// Find @response sequence with shorthand (Target != "").
	var varName string
	for _, seq := range fn.Sequences {
		if seq.Type == "response" && seq.Target != "" {
			varName = seq.Target
			break
		}
	}
	if varName == "" {
		return nil
	}

	// Trace variable origin: find the sequence that assigned this variable.
	for _, seq := range fn.Sequences {
		if seq.Result == nil || seq.Result.Var != varName {
			continue
		}

		// Wrapper types (Page, Cursor) have fixed structure — skip.
		if seq.Result.Wrapper != "" {
			return nil
		}

		switch seq.Type {
		case "call":
			// @call — resolve from funcspec ResponseFields.
			// Result.Type may include package prefix (e.g., "auth.IssueTokenResponse").
			typeName := seq.Result.Type
			if idx := strings.LastIndex(typeName, "."); idx >= 0 {
				typeName = typeName[idx+1:]
			}
			for _, fs := range funcSpecs {
				expectedResp := ucFirst(fs.Name) + "Response"
				if expectedResp != typeName {
					continue
				}
				var keys []string
				for _, f := range fs.ResponseFields {
					if f.JSONName != "" {
						keys = append(keys, f.JSONName)
					} else {
						keys = append(keys, f.Name)
					}
				}
				return keys
			}

		case "get", "put", "post", "delete":
			// @get/@put/@post/@delete — resolve from DDL columns.
			if st == nil {
				return nil
			}
			tableName := seq.Result.Type
			// DDL table names are lowercase plural (e.g., Type="Gig" → table="gigs").
			// Try symbol table lookup with various conventions.
			for tbl, ddl := range st.DDLTables {
				// Match: type name matches table name (case-insensitive singular).
				if matchTableType(tbl, tableName) {
					var keys []string
					for col := range ddl.Columns {
						keys = append(keys, col)
					}
					sort.Strings(keys)
					return keys
				}
			}
		}
		break
	}

	return nil
}

// matchTableType checks if a DDL table name corresponds to a type name.
// e.g., "gigs" matches "Gig", "users" matches "User".
func matchTableType(tableName, typeName string) bool {
	tn := strings.ToLower(typeName)
	tbl := strings.ToLower(tableName)
	// Exact match or simple plural (type + "s" == table).
	return tbl == tn || tbl == tn+"s" || tbl == tn+"es"
}

// ucFirst converts first letter to uppercase (PascalCase).
func ucFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// extractResponseFieldKeys returns the @response field keys for a function,
// or nil if the function uses shorthand (@response varName) or has no @response.
func extractResponseFieldKeys(fn ssacparser.ServiceFunc) []string {
	for _, seq := range fn.Sequences {
		if seq.Type != "response" {
			continue
		}
		// Shorthand: @response varName — no individual field keys.
		if seq.Target != "" {
			return nil
		}
		if len(seq.Fields) == 0 {
			return nil
		}
		keys := make([]string, 0, len(seq.Fields))
		for k := range seq.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}
	return nil
}

// buildOperationResponseProps collects response schema property names per operationId from the OpenAPI doc.
func buildOperationResponseProps(doc *openapi3.T) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	if doc == nil || doc.Paths == nil {
		return result
	}

	for _, pathItem := range doc.Paths.Map() {
		for _, op := range []*openapi3.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put,
			pathItem.Delete, pathItem.Patch,
		} {
			if op == nil || op.OperationID == "" || op.Responses == nil {
				continue
			}

			props := make(map[string]bool)

			// Check all 2xx responses.
			for code, respRef := range op.Responses.Map() {
				if len(code) != 3 || code[0] != '2' {
					continue
				}
				if respRef == nil || respRef.Value == nil || respRef.Value.Content == nil {
					continue
				}
				ct := respRef.Value.Content.Get("application/json")
				if ct == nil || ct.Schema == nil {
					continue
				}
				schema := resolveSchemaRef(ct.Schema)
				if schema == nil {
					continue
				}
				for propName := range schema.Properties {
					props[propName] = true
				}
			}

			if len(props) > 0 {
				result[op.OperationID] = props
			}
		}
	}

	return result
}

func resolveSchemaRef(ref *openapi3.SchemaRef) *openapi3.Schema {
	if ref == nil {
		return nil
	}
	return ref.Value
}

// checkResponseSuccessCode validates that SSaC functions with @response have an explicit 2xx response code in OpenAPI.
// "default"-only responses are not accepted — the OpenAPI spec must explicitly declare 200, 201, 204, etc.
func checkResponseSuccessCode(funcs []ssacparser.ServiceFunc, doc *openapi3.T) []CrossError {
	var errs []CrossError

	// Build operationId → operation map.
	opMap := make(map[string]*openapi3.Operation)
	if doc.Paths != nil {
		for _, pathItem := range doc.Paths.Map() {
			for _, op := range []*openapi3.Operation{
				pathItem.Get, pathItem.Post, pathItem.Put,
				pathItem.Delete, pathItem.Patch,
			} {
				if op != nil && op.OperationID != "" {
					opMap[op.OperationID] = op
				}
			}
		}
	}

	for _, fn := range funcs {
		hasResponse := false
		for _, seq := range fn.Sequences {
			if seq.Type == "response" {
				hasResponse = true
				break
			}
		}
		if !hasResponse {
			continue
		}

		op := opMap[fn.Name]
		if op == nil || op.Responses == nil {
			continue // no matching operation — already caught by Rule 3
		}

		// Check for explicit 2xx response code.
		has2xx := false
		for code := range op.Responses.Map() {
			if len(code) == 3 && code[0] == '2' {
				has2xx = true
				break
			}
		}

		if !has2xx {
			errs = append(errs, CrossError{
				Rule:       "SSaC @response → OpenAPI 2xx",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("SSaC @response가 있는 %s에 OpenAPI 2xx 성공 응답 코드가 없습니다 (default만으로는 불충분)", fn.Name),
				Suggestion: fmt.Sprintf("OpenAPI %s responses에 200, 201, 204 등 명시적 성공 코드를 추가하세요", fn.Name),
			})
		}
	}

	return errs
}

// errStatusTypes are the SSaC sequence types that support custom ErrStatus.
var errStatusTypes = map[string]int{
	"empty": 404,
	"exists": 409,
	"state": 409,
	"auth":  403,
}

// checkErrStatus validates that SSaC ErrStatus codes are defined in OpenAPI responses.
func checkErrStatus(funcs []ssacparser.ServiceFunc, doc *openapi3.T) []CrossError {
	var errs []CrossError

	// Build operationId → operation map.
	opMap := make(map[string]*openapi3.Operation)
	if doc.Paths != nil {
		for _, pathItem := range doc.Paths.Map() {
			for _, op := range []*openapi3.Operation{
				pathItem.Get, pathItem.Post, pathItem.Put,
				pathItem.Delete, pathItem.Patch,
			} {
				if op != nil && op.OperationID != "" {
					opMap[op.OperationID] = op
				}
			}
		}
	}

	for _, fn := range funcs {
		op := opMap[fn.Name]
		if op == nil || op.Responses == nil {
			continue
		}

		for seqIdx, seq := range fn.Sequences {
			defaultStatus, ok := errStatusTypes[seq.Type]
			if !ok {
				continue
			}

			statusCode := defaultStatus
			if seq.ErrStatus != 0 {
				statusCode = seq.ErrStatus
			}

			codeStr := fmt.Sprintf("%d", statusCode)
			resp := op.Responses.Status(statusCode)
			if resp == nil {
				errs = append(errs, CrossError{
					Rule:       "SSaC @" + seq.Type + " → OpenAPI",
					Context:    fmt.Sprintf("%s:%s seq[%d]", fn.FileName, fn.Name, seqIdx),
					Message:    fmt.Sprintf("SSaC @%s uses HTTP %s but OpenAPI %s has no %s response defined", seq.Type, codeStr, fn.Name, codeStr),
					Suggestion: fmt.Sprintf("OpenAPI %s responses에 %s 응답을 추가하세요", fn.Name, codeStr),
				})
			}
		}
	}

	return errs
}
