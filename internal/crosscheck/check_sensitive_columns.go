//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=sensitive
//ff:what DDL 컬럼 이름이 민감 패턴에 매치되지만 @sensitive 없는 경우 경고
package crosscheck

import (
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// sensitivePatterns are column name substrings that suggest sensitive data.
var sensitivePatterns = []string{
	// 인증 정보
	"password", "passwd", "passphrase",
	"secret", "token", "hash", "salt",
	"credential", "otp", "pin",
	// 암호화
	"private_key", "cipher", "encrypted",
	// 금융
	"credit_card", "card_number", "cvv",
	"bank_account", "routing_number",
	// 개인식별
	"ssn", "passport", "license_number",
	"biometric",
}

// CheckSensitiveColumns warns when DDL column names match sensitive patterns
// but lack an @sensitive annotation.
func CheckSensitiveColumns(st *ssacvalidator.SymbolTable, sensitiveCols, noSensitiveCols map[string]map[string]bool) []CrossError {
	var errs []CrossError

	for tableName, table := range st.DDLTables {
		errs = append(errs, checkTableSensitiveColumns(tableName, table.ColumnOrder, sensitiveCols, noSensitiveCols)...)
	}

	return errs
}
