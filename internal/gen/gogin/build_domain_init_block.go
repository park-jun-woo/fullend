//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=output
//ff:what 도메인 핸들러 초기화 코드 블록을 생성한다

package gogin

import (
	"fmt"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func buildDomainInitBlock(serviceFuncs []ssacparser.ServiceFunc, domains []string, anyNeedsAuth bool) string {
	flatModels := collectModelsForDomain(serviceFuncs, "")

	var initLines []string

	for _, m := range flatModels {
		fieldName := ucFirst(lcFirst(m) + "Model")
		initLines = append(initLines, fmt.Sprintf("\t\t%s: model.New%sModel(conn),", fieldName, m))
	}

	for _, domain := range domains {
		domainModels := collectModelsForDomain(serviceFuncs, domain)
		fieldName := ucFirst(domain)

		var handlerLines []string
		if domainNeedsDB(serviceFuncs, domain) {
			handlerLines = append(handlerLines, "\t\t\tDB: conn,")
		}
		for _, m := range domainModels {
			mFieldName := ucFirst(lcFirst(m) + "Model")
			handlerLines = append(handlerLines, fmt.Sprintf("\t\t\t%s: model.New%sModel(conn),", mFieldName, m))
		}
		if domainNeedsJWTSecret(serviceFuncs, domain) {
			handlerLines = append(handlerLines, "\t\t\tJWTSecret: *jwtSecret,")
		}
		initLines = append(initLines, fmt.Sprintf("\t\t%s: &%ssvc.Handler{", fieldName, domain))
		initLines = append(initLines, handlerLines...)
		initLines = append(initLines, "\t\t},")
	}

	if anyNeedsAuth {
		initLines = append(initLines, "\t\tJWTSecret: *jwtSecret,")
	}

	initBlock := strings.Join(initLines, "\n")
	if initBlock == "" {
		initBlock = "\t\t// No models detected"
	}
	return initBlock
}
