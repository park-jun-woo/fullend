//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what 교차 검증 실행 — OpenAPI+DDL+SSaC 교차 정합성 검증 오케스트레이션
package orchestrator

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/crosscheck"
	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func runCrossValidate(root string, parsed *genapi.ParsedSSOTs) reporter.StepResult {
	step := reporter.StepResult{Name: "Cross"}

	// Require OpenAPI + DDL + SSaC for cross-validation.
	if parsed.OpenAPIDoc == nil || parsed.SymbolTable == nil || parsed.ServiceFuncs == nil {
		step.Status = reporter.Skip
		step.Summary = "skipped (incomplete SSOT)"
		return step
	}

	// Load @dto types from model files.
	dtoTypes := loadDTOTypes(parsed.ModelDir)

	var middleware []string
	var claims map[string]projectconfig.ClaimDef
	var roles []string
	if parsed.Config != nil {
		middleware = parsed.Config.Backend.Middleware
		if parsed.Config.Backend.Auth != nil {
			claims = parsed.Config.Backend.Auth.Claims
			roles = parsed.Config.Backend.Auth.Roles
		}
	}

	// Parse @archived tags from DDL files.
	archived, _ := crosscheck.ParseArchived(filepath.Join(root, "db"))

	// Parse @sensitive / @nosensitive tags from DDL files.
	sensitiveCols, noSensitiveCols, _ := crosscheck.ParseSensitive(filepath.Join(root, "db"))

	var queueBackend string
	if parsed.Config != nil && parsed.Config.Queue != nil {
		queueBackend = parsed.Config.Queue.Backend
	}

	var authzPackage string
	if parsed.Config != nil && parsed.Config.Authz != nil {
		authzPackage = parsed.Config.Authz.Package
	}

	input := &crosscheck.CrossValidateInput{
		ParsedSSOTs:     parsed,
		DTOTypes:        dtoTypes,
		Middleware:      middleware,
		Archived:        archived,
		Claims:          claims,
		QueueBackend:    queueBackend,
		AuthzPackage:    authzPackage,
		SensitiveCols:   sensitiveCols,
		NoSensitiveCols: noSensitiveCols,
		Roles:           roles,
	}

	cerrs := crosscheck.Run(input)

	hasError := false
	for _, ce := range cerrs {
		prefix := ce.Rule
		if ce.Level == "WARNING" {
			prefix = "[WARN] " + prefix
		} else {
			hasError = true
		}
		step.Errors = append(step.Errors, fmt.Sprintf("%s: %s — %s", prefix, ce.Context, ce.Message))
		step.Suggestions = append(step.Suggestions, ce.Suggestion)
	}

	if hasError {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}

	errCount := 0
	warnCount := 0
	for _, ce := range cerrs {
		if ce.Level == "WARNING" {
			warnCount++
		} else {
			errCount++
		}
	}
	if errCount > 0 {
		step.Summary = fmt.Sprintf("%d errors, %d warnings", errCount, warnCount)
	} else if warnCount > 0 {
		step.Summary = fmt.Sprintf("%d warnings", warnCount)
	} else {
		step.Summary = "0 mismatches"
	}
	return step
}
