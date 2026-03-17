package crosscheck

import (
	"testing"

	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckSensitiveColumns_NewPatterns(t *testing.T) {
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {
				ColumnOrder: []string{"id", "user_credential", "otp_code", "api_key", "ssn_number"},
				Columns: map[string]string{
					"id":              "int64",
					"user_credential": "string",
					"otp_code":        "string",
					"api_key":         "string",
					"ssn_number":      "string",
				},
			},
		},
	}

	errs := CheckSensitiveColumns(st, nil, nil)

	matched := make(map[string]bool)
	for _, e := range errs {
		matched[e.Context] = true
	}

	// Should detect credential, otp, ssn
	for _, want := range []string{"users.user_credential", "users.otp_code", "users.ssn_number"} {
		if !matched[want] {
			t.Errorf("expected sensitive warning for %s", want)
		}
	}

	// api_key should NOT be detected (key excluded from patterns)
	if matched["users.api_key"] {
		t.Error("api_key should not be detected (key excluded from patterns)")
	}
}
