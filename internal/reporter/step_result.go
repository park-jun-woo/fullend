//ff:type feature=reporter type=model
//ff:what 단일 검증 단계의 결과를 담는 구조체
package reporter

// StepResult holds the outcome of a single validation step.
type StepResult struct {
	Name        string   // "OpenAPI", "DDL", "SSaC", "STML", "Cross"
	Status      Status   // Pass, Fail, Skip
	Summary     string   // "34 endpoints", "12 tables, 47 columns"
	Errors      []string // individual error messages
	Suggestions []string // fix suggestions (parallel to Errors, empty string if none)
}
