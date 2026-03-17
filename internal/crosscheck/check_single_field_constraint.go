//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what 단일 필드의 C2~C4 ERROR, W1~W3 WARNING 규칙을 적용
package crosscheck

import (
	"fmt"
	"strings"

	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkSingleFieldConstraint(opID, fieldName, col string, fc ssacvalidator.FieldConstraint, varcharLen int, checkEnums []string, found bool) []CrossError {
	var errs []CrossError

	if found && varcharLen > 0 && fc.MaxLength == nil {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints C2",
			Context:    opID,
			Message:    fmt.Sprintf("DDL column %q is VARCHAR(%d) but OpenAPI field %q has no maxLength", col, varcharLen, fieldName),
			Suggestion: fmt.Sprintf("OpenAPI에 maxLength: %d 추가", varcharLen),
		})
	}

	if found && len(checkEnums) > 0 && len(fc.Enum) == 0 {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints C3",
			Context:    opID,
			Message:    fmt.Sprintf("DDL column %q has CHECK IN (%s) but OpenAPI field %q has no enum", col, strings.Join(checkEnums, ", "), fieldName),
			Suggestion: fmt.Sprintf("OpenAPI에 enum: [%s] 추가", strings.Join(checkEnums, ", ")),
		})
	}

	if found && len(checkEnums) > 0 && len(fc.Enum) > 0 && !enumsMatch(checkEnums, fc.Enum) {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints C4",
			Context:    opID,
			Message:    fmt.Sprintf("DDL CHECK IN (%s) != OpenAPI enum [%s] for field %q", strings.Join(checkEnums, ", "), strings.Join(fc.Enum, ", "), fieldName),
			Suggestion: "DDL CHECK 값과 OpenAPI enum 값을 일치시키세요",
		})
	}

	if found && varcharLen > 0 && fc.MaxLength != nil && *fc.MaxLength > varcharLen {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints W1",
			Context:    opID,
			Level:      "WARNING",
			Message:    fmt.Sprintf("OpenAPI maxLength(%d) > DDL VARCHAR(%d) for field %q", *fc.MaxLength, varcharLen, fieldName),
			Suggestion: fmt.Sprintf("maxLength를 %d 이하로 조정", varcharLen),
		})
	}

	if isPasswordField(fieldName) && fc.MinLength == nil {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints W2",
			Context:    opID,
			Level:      "WARNING",
			Message:    fmt.Sprintf("password field %q has no minLength constraint", fieldName),
			Suggestion: "보안을 위해 minLength 제약 추가 (예: minLength: 8)",
		})
	}

	if isEmailField(fieldName) && fc.Format != "email" {
		errs = append(errs, CrossError{
			Rule:       "OpenAPI Constraints W3",
			Context:    opID,
			Level:      "WARNING",
			Message:    fmt.Sprintf("email field %q has no format: email constraint", fieldName),
			Suggestion: "format: email 추가",
		})
	}

	return errs
}
