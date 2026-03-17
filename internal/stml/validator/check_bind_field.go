//ff:func feature=stml-validate type=rule control=sequence
//ff:what 단일 data-bind 필드가 응답 스키마 또는 custom.ts에 있는지 확인
package validator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

func checkBindField(b parser.FieldBind, opID, file string, api APISymbol, cs *CustomSymbol) *ValidationError {
	fieldName := b.Name
	if idx := strings.IndexByte(fieldName, '.'); idx >= 0 {
		fieldName = fieldName[:idx]
	}
	if _, ok := api.ResponseFields[fieldName]; ok {
		return nil
	}
	if cs != nil && cs.Functions[fieldName] {
		return nil
	}
	err := errBindNotFound(file, opID, b.Name)
	return &err
}
