//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what Generate — Fullstack + Ground 에서 Go+Gin 백엔드 코드 생성
package gogin

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Config holds Go+Gin backend generation configuration.
type Config struct {
	ArtifactsDir string
	SpecsDir     string
	ModulePath   string
}

// Generate creates Go+Gin backend code from Fullstack + Ground.
func (g *GoGin) Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
	intDir := internalDir(cfg.ArtifactsDir)
	if err := os.MkdirAll(intDir, 0755); err != nil {
		return fmt.Errorf("create internal dir: %w", err)
	}

	// Config 추출
	var claims map[string]manifest.ClaimDef
	var secretEnv string
	var queueBackend string
	var sessionBackend, cacheBackend string
	var fileConfig *manifest.FileBackend
	if fs.Manifest != nil {
		if fs.Manifest.Backend.Auth != nil {
			claims = fs.Manifest.Backend.Auth.Claims
			secretEnv = fs.Manifest.Backend.Auth.SecretEnv
		}
		if fs.Manifest.Queue != nil {
			queueBackend = fs.Manifest.Queue.Backend
		}
		if fs.Manifest.Session != nil {
			sessionBackend = fs.Manifest.Session.Backend
		}
		if fs.Manifest.Cache != nil {
			cacheBackend = fs.Manifest.Cache.Backend
		}
		fileConfig = fs.Manifest.File
	}

	if hasBearerScheme(fs.OpenAPIDoc) && len(claims) == 0 {
		return fmt.Errorf("OpenAPI has bearerAuth security scheme but fullend.yaml has no claims config")
	}

	models := collectModels(fs.ServiceFuncs)
	funcs := collectFuncs(fs.ServiceFuncs)
	policies := adaptPolicies(fs.ParsedPolicies)

	// Feature mode (Flat 제거됨)
	if err := transformServiceFiles(intDir, fs.ServiceFuncs, models, funcs, cfg.ModulePath, fs.OpenAPIDoc); err != nil {
		return fmt.Errorf("service transform: %w", err)
	}
	if err := attachServiceDirectives(intDir, fs.ServiceFuncs); err != nil {
		return fmt.Errorf("service directives: %w", err)
	}
	if err := generateAuthStub(intDir, cfg.ModulePath, claims); err != nil {
		return fmt.Errorf("auth stub: %w", err)
	}
	if err := generateAuthIfNeeded(intDir, cfg.ModulePath, claims, secretEnv); err != nil {
		return err
	}
	if err := generateServerStruct(intDir, fs.ServiceFuncs, cfg.ModulePath, fs.OpenAPIDoc); err != nil {
		return fmt.Errorf("server.go: %w", err)
	}
	if err := generateMain(MainGenInput{
		ArtifactsDir:   cfg.ArtifactsDir,
		ServiceFuncs:   fs.ServiceFuncs,
		ModulePath:     cfg.ModulePath,
		QueueBackend:   queueBackend,
		Policies:       policies,
		SessionBackend: sessionBackend,
		CacheBackend:   cacheBackend,
		FileConfig:     fileConfig,
	}); err != nil {
		return fmt.Errorf("main.go: %w", err)
	}

	// 공유
	modelIncludeSpecs := collectModelIncludes(fs.OpenAPIDoc, fs.ServiceFuncs)
	cursorSpecs := collectCursorSpecs(fs.OpenAPIDoc)
	if err := generateModelImpls(intDir, models, cfg.ModulePath, cfg.SpecsDir, fs.ServiceFuncs, modelIncludeSpecs, cursorSpecs); err != nil {
		return fmt.Errorf("model impl: %w", err)
	}
	if err := attachTSXDirectives(cfg.ArtifactsDir); err != nil {
		return fmt.Errorf("tsx directives: %w", err)
	}

	return nil
}
