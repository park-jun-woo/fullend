//ff:func feature=rule type=generator control=iteration dimension=1
//ff:what generateBackend — 각 ServiceFunc의 핸들러 body 생성
package generate

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/pkg/generate/backend"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func generateBackend(fs *fullend.Fullstack, artifactsDir string) []string {
	serviceDir := filepath.Join(artifactsDir, "backend", "internal", "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return []string{"cannot create service dir: " + err.Error()}
	}
	var errs []string
	for _, fn := range fs.ServiceFuncs {
		if err := writeHandler(serviceDir, fn, backend.GenerateHandler(fn)); err != "" {
			errs = append(errs, err)
		}
	}
	return errs
}
