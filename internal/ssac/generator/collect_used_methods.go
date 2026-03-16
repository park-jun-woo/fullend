//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=model-collect
//ff:what 특정 모델에서 사용된 메서드명을 중복 없이 수집
package generator

func collectUsedMethods(modelName string, usages []modelUsage) []string {
	var usedMethods []string
	for _, u := range usages {
		if u.ModelName != modelName {
			continue
		}
		if !containsString(usedMethods, u.MethodName) {
			usedMethods = append(usedMethods, u.MethodName)
		}
	}
	return usedMethods
}
