//ff:func feature=gen-gogin type=generator control=selection
//ff:what fileConfig → import/init 코드 조각 반환 (local/s3 분기)

package gogin

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

func buildFileInitBlock(fileConfig *manifest.FileBackend) (imports []string, initSnippet string) {
	imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/file"`)

	switch fileConfig.Backend {
	case "local":
		root := "./uploads"
		if fileConfig.Local != nil && fileConfig.Local.Root != "" {
			root = fileConfig.Local.Root
		}
		initSnippet = fmt.Sprintf(`
	file.Init(file.NewLocalFile(%q))`, root)
	case "s3":
		bucket := ""
		region := "ap-northeast-2"
		if fileConfig.S3 != nil {
			bucket = fileConfig.S3.Bucket
			region = fileConfig.S3.Region
		}
		imports = append(imports, `"github.com/aws/aws-sdk-go-v2/config"`)
		imports = append(imports, `"github.com/aws/aws-sdk-go-v2/service/s3"`)
		initSnippet = fmt.Sprintf(`
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(%q))
	if err != nil {
		log.Fatalf("aws config failed: %%v", err)
	}
	file.Init(file.NewS3File(s3.NewFromConfig(awsCfg), %q))`, region, bucket)
	}
	return imports, initSnippet
}
