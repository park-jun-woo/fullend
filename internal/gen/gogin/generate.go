//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what main entry point — orchestrates Go+Gin backend code generation from parsed SSOTs

package gogin

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

// Generate creates Go+Gin backend code from parsed SSOTs.
func (g *GoGin) Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig) error {
	intDir := internalDir(cfg.ArtifactsDir)
	if err := os.MkdirAll(intDir, 0755); err != nil {
		return fmt.Errorf("create internal dir: %w", err)
	}

	// Extract config values.
	var claims map[string]projectconfig.ClaimDef
	var secretEnv string
	var queueBackend string
	var sessionBackend, cacheBackend string
	var fileConfig *projectconfig.FileBackend
	if parsed.Config != nil {
		if parsed.Config.Backend.Auth != nil {
			claims = parsed.Config.Backend.Auth.Claims
			secretEnv = parsed.Config.Backend.Auth.SecretEnv
		}
		if parsed.Config.Queue != nil {
			queueBackend = parsed.Config.Queue.Backend
		}
		if parsed.Config.Session != nil {
			sessionBackend = parsed.Config.Session.Backend
		}
		if parsed.Config.Cache != nil {
			cacheBackend = parsed.Config.Cache.Backend
		}
		fileConfig = parsed.Config.File
	}

	// Validate: bearerAuth requires claims config.
	if hasBearerScheme(parsed.OpenAPIDoc) && len(claims) == 0 {
		return fmt.Errorf("OpenAPI has bearerAuth security scheme but fullend.yaml has no claims config")
	}

	models := collectModels(parsed.ServiceFuncs)

	if hasDomains(parsed.ServiceFuncs) {
		// Domain mode: per-domain Handler + central Server.
		allFuncs := collectFuncs(parsed.ServiceFuncs)

		if err := transformServiceFilesWithDomains(intDir, parsed.ServiceFuncs, models, allFuncs, cfg.ModulePath, parsed.OpenAPIDoc); err != nil {
			return fmt.Errorf("service transform (domain): %w", err)
		}
		if err := attachServiceDirectives(intDir, parsed.ServiceFuncs); err != nil {
			return fmt.Errorf("service directives (domain): %w", err)
		}
		if err := generateAuthStubWithDomains(intDir, cfg.ModulePath, claims); err != nil {
			return fmt.Errorf("auth (domain): %w", err)
		}
		if err := generateAuthIfNeeded(intDir, cfg.ModulePath, claims, secretEnv); err != nil {
			return err
		}
		if err := generateServerStructWithDomains(intDir, parsed.ServiceFuncs, cfg.ModulePath, parsed.OpenAPIDoc); err != nil {
			return fmt.Errorf("server.go (domain): %w", err)
		}
		if err := generateMainWithDomains(cfg.ArtifactsDir, parsed.ServiceFuncs, cfg.ModulePath, queueBackend, parsed.Policies, sessionBackend, cacheBackend, fileConfig); err != nil {
			return fmt.Errorf("main.go (domain): %w", err)
		}
	} else {
		// Flat mode: single Server with all fields (unchanged).
		funcs := collectFuncs(parsed.ServiceFuncs)

		if err := transformServiceFiles(intDir, models, funcs, cfg.ModulePath, parsed.OpenAPIDoc, parsed.ServiceFuncs); err != nil {
			return fmt.Errorf("service transform: %w", err)
		}
		if err := attachServiceDirectives(intDir, parsed.ServiceFuncs); err != nil {
			return fmt.Errorf("service directives: %w", err)
		}
		if err := generateServerStruct(intDir, models, funcs, cfg.ModulePath, parsed.OpenAPIDoc); err != nil {
			return fmt.Errorf("server.go: %w", err)
		}
		if err := generateAuthStub(intDir); err != nil {
			return fmt.Errorf("auth.go: %w", err)
		}
		if err := generateMain(cfg.ArtifactsDir, models, cfg.ModulePath, queueBackend, parsed.ServiceFuncs, parsed.Policies, sessionBackend, cacheBackend, fileConfig); err != nil {
			return fmt.Errorf("main.go: %w", err)
		}
	}

	// Shared: model implementations + TSX directives (same for both modes).
	modelIncludeSpecs := collectModelIncludes(parsed.OpenAPIDoc, parsed.ServiceFuncs)
	cursorSpecs := collectCursorSpecs(parsed.OpenAPIDoc)
	if err := generateModelImpls(intDir, models, cfg.ModulePath, cfg.SpecsDir, parsed.ServiceFuncs, modelIncludeSpecs, cursorSpecs); err != nil {
		return fmt.Errorf("model impl: %w", err)
	}
	if err := attachTSXDirectives(cfg.ArtifactsDir); err != nil {
		return fmt.Errorf("tsx directives: %w", err)
	}

	return nil
}
