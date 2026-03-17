//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what session/cache/file Init 코드 블록 + import 빌더

package gogin

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func buildBuiltinInitBlocks(sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend) (builtinImport, builtinInitBlock string) {
	var imports []string
	var inits []string

	// --- session ---
	if sessionBackend == "postgres" {
		imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/session"`)
		inits = append(inits, `
	sm, err := session.NewPostgresSession(context.Background(), conn)
	if err != nil {
		log.Fatalf("session init failed: %v", err)
	}
	session.Init(sm)`)
	} else if sessionBackend == "memory" {
		imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/session"`)
		inits = append(inits, `
	session.Init(session.NewMemorySession())`)
	}

	// --- cache ---
	if cacheBackend == "postgres" {
		imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/cache"`)
		inits = append(inits, `
	cm, err := cache.NewPostgresCache(context.Background(), conn)
	if err != nil {
		log.Fatalf("cache init failed: %v", err)
	}
	cache.Init(cm)`)
	} else if cacheBackend == "memory" {
		imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/cache"`)
		inits = append(inits, `
	cache.Init(cache.NewMemoryCache())`)
	}

	// --- file ---
	if fileConfig != nil {
		imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/file"`)
		switch fileConfig.Backend {
		case "local":
			root := "./uploads"
			if fileConfig.Local != nil && fileConfig.Local.Root != "" {
				root = fileConfig.Local.Root
			}
			inits = append(inits, fmt.Sprintf(`
	file.Init(file.NewLocalFile(%q))`, root))
		case "s3":
			bucket := ""
			region := "ap-northeast-2"
			if fileConfig.S3 != nil {
				bucket = fileConfig.S3.Bucket
				region = fileConfig.S3.Region
			}
			imports = append(imports, `"github.com/aws/aws-sdk-go-v2/config"`)
			imports = append(imports, `"github.com/aws/aws-sdk-go-v2/service/s3"`)
			inits = append(inits, fmt.Sprintf(`
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(%q))
	if err != nil {
		log.Fatalf("aws config failed: %%v", err)
	}
	file.Init(file.NewS3File(s3.NewFromConfig(awsCfg), %q))`, region, bucket))
		}
	}

	// --- context import (postgres session/cache 또는 s3에서 사용) ---
	needsContext := sessionBackend == "postgres" || cacheBackend == "postgres" ||
		(fileConfig != nil && fileConfig.Backend == "s3")
	if needsContext {
		imports = append([]string{`"context"`}, imports...)
	}

	if len(imports) > 0 {
		builtinImport = "\n\t" + strings.Join(imports, "\n\t")
	}
	if len(inits) > 0 {
		builtinInitBlock = strings.Join(inits, "")
	}
	return builtinImport, builtinInitBlock
}
