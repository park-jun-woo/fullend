//ff:func feature=gen-gogin type=generator control=sequence
//ff:what SSaC 생성 소스의 import 경로를 프로젝트 모듈 경로로 변환한다

package gogin

import (
	"fmt"
	"strings"
)

// fixImports rewrites short import paths in SSaC-generated source to full module paths.
func fixImports(src string, modulePath string) string {
	// Add model import when model package types are used (QueryOpts, ParseQueryOpts, CurrentUser, etc.)
	needsModel := strings.Contains(src, "model.QueryOpts") || strings.Contains(src, "model.ParseQueryOpts") || strings.Contains(src, "*model.CurrentUser")
	if needsModel {
		modelImport := fmt.Sprintf("\"%s/internal/model\"", modulePath)
		if !strings.Contains(src, modelImport) {
			src = strings.Replace(src, "import (", "import (\n\t"+modelImport, 1)
		}
	}

	// Fix state machine imports: "states/XXXstate" -> "{modulePath}/internal/states/XXXstate"
	if strings.Contains(src, "\"states/") {
		src = strings.ReplaceAll(src, "\"states/", fmt.Sprintf("\"%s/internal/states/", modulePath))
	}

	// Fix authz import: "authz" -> "github.com/park-jun-woo/fullend/pkg/authz"
	if strings.Contains(src, "\t\"authz\"\n") {
		src = strings.ReplaceAll(src, "\t\"authz\"\n", "\t\"github.com/park-jun-woo/fullend/pkg/authz\"\n")
	}

	// Fix queue import: "queue" -> "github.com/park-jun-woo/fullend/pkg/queue"
	if strings.Contains(src, "\t\"queue\"\n") {
		src = strings.ReplaceAll(src, "\t\"queue\"\n", "\t\"github.com/park-jun-woo/fullend/pkg/queue\"\n")
	}

	// Fix config import: "config" -> "github.com/park-jun-woo/fullend/pkg/config"
	if strings.Contains(src, "\t\"config\"\n") {
		src = strings.ReplaceAll(src, "\t\"config\"\n", "\t\"github.com/park-jun-woo/fullend/pkg/config\"\n")
	}

	// Fix auth import: pkg/auth -> project internal/auth (reexport.go bridges pkg/auth utilities)
	if strings.Contains(src, "\"github.com/park-jun-woo/fullend/pkg/auth\"") {
		src = strings.ReplaceAll(src, "\"github.com/park-jun-woo/fullend/pkg/auth\"",
			fmt.Sprintf("\"%s/internal/auth\"", modulePath))
	}

	// Remove bare "model" import (already added as full path above)
	if strings.Contains(src, "\t\"model\"\n") {
		src = strings.ReplaceAll(src, "\t\"model\"\n", "")
	}

	return src
}
