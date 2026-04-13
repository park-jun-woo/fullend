//ff:func feature=rule type=util control=sequence
//ff:what ensureModel — g.Models[key] 이 없으면 빈 엔트리 생성 후 반환
package ground

import "github.com/park-jun-woo/fullend/pkg/rule"

func ensureModel(g *rule.Ground, key string) rule.ModelInfo {
	info, ok := g.Models[key]
	if !ok {
		info = rule.ModelInfo{Name: key, Methods: make(map[string]rule.MethodInfo)}
	}
	if info.Methods == nil {
		info.Methods = make(map[string]rule.MethodInfo)
	}
	return info
}
