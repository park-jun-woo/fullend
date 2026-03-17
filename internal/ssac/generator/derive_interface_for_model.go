//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=interface-derive
//ff:what 단일 모델의 메서드 사용을 분석하여 인터페이스를 파생
package generator

import (
	"sort"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func deriveInterfaceForModel(modelName string, usages []modelUsage, methodMap map[methodKey]modelUsage, st *validator.SymbolTable) *derivedInterface {
	ms, ok := st.Models[modelName]
	if !ok {
		return nil
	}

	iface := derivedInterface{Name: modelName + "Model"}
	usedMethods := collectUsedMethods(modelName, usages)
	sort.Strings(usedMethods)

	for _, methodName := range usedMethods {
		mi, methodExists := ms.Methods[methodName]
		if !methodExists {
			mi = validator.MethodInfo{}
		}
		key := methodKey{modelName, methodName}
		usage := methodMap[key]

		dm := deriveMethod(methodName, usage, mi, st)
		iface.Methods = append(iface.Methods, dm)
	}

	if len(iface.Methods) == 0 {
		return nil
	}
	return &iface
}
