//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=output
//ff:what 도메인 import 코드 블록을 생성한다

package gogin

import (
	"fmt"
	"strings"
)

func buildDomainImportsBlock(domains []string, modulePath string, anyNeedsAuth bool) string {
	var extraImports []string
	extraImports = append(extraImports, fmt.Sprintf("\n\t\"%s/internal/model\"", modulePath))
	extraImports = append(extraImports, fmt.Sprintf("\t\"%s/internal/service\"", modulePath))
	if anyNeedsAuth {
		extraImports = append(extraImports, "\t\"github.com/park-jun-woo/fullend/pkg/authz\"")
	}
	for _, d := range domains {
		extraImports = append(extraImports, fmt.Sprintf("\t%ssvc \"%s/internal/service/%s\"", d, modulePath, d))
	}
	return strings.Join(extraImports, "\n")
}
