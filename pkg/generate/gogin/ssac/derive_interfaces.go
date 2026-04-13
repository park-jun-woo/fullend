//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=interface-derive
//ff:what 모델 사용 정보와 심볼 테이블에서 인터페이스를 파생
package ssac

import "github.com/park-jun-woo/fullend/internal/ssac/validator"

func deriveInterfaces(usages []modelUsage, st *validator.SymbolTable) []derivedInterface {
	methodMap := map[methodKey]modelUsage{}
	modelNames := map[string]bool{}

	for _, u := range usages {
		key := methodKey{u.ModelName, u.MethodName}
		if _, exists := methodMap[key]; !exists {
			methodMap[key] = u
			modelNames[u.ModelName] = true
		}
	}

	var interfaces []derivedInterface
	sortedModels := sortedKeys(modelNames)

	for _, modelName := range sortedModels {
		iface := deriveInterfaceForModel(modelName, usages, methodMap, st)
		if iface != nil {
			interfaces = append(interfaces, *iface)
		}
	}

	return interfaces
}
