//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=interface-derive
//ff:what 단일 메서드의 파라미터와 반환 타입을 파생
package generator

import (
	"github.com/ettle/strcase"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func deriveMethod(methodName string, usage modelUsage, mi validator.MethodInfo, st *validator.SymbolTable) derivedMethod {
	dm := derivedMethod{Name: methodName}

	inputKeys := orderMethodInputKeys(usage.Inputs, mi.Params)

	for _, k := range inputKeys {
		val := usage.Inputs[k]
		if val == "query" {
			dm.HasQueryOpts = true
			continue
		}
		goType := resolveParamTypeWithFallback(val, k, usage.ModelName, st)
		dp := derivedParam{
			Name:   strcase.ToGoCamel(k),
			GoType: goType,
		}
		if dp.Name != "" {
			dm.Params = append(dm.Params, dp)
		}
	}

	dm.ReturnType = deriveReturnType(mi, usage, dm.HasQueryOpts)
	return dm
}
