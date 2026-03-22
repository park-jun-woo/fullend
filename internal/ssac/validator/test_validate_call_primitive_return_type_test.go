//ff:func feature=ssac-validate type=test control=iteration dimension=1
//ff:what @call 결과가 기본 타입이면 ERROR 검증
package validator

import (
	"testing"
)

func TestValidateCallPrimitiveReturnType(t *testing.T) {
	for _, typ := range []string{"string", "int", "int64", "bool", "float64", "time.Time"} {
		t.Run(typ, func(t *testing.T) {
			assertCallPrimitiveCase(t, typ)
		})
	}
}
