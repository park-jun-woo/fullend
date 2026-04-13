//ff:func feature=rule type=loader control=iteration dimension=2
//ff:what seedFromInterfaces — iface.Interface 배열로 g.Models 초기 메서드 맵 구축
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/parser/iface"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func seedFromInterfaces(g *rule.Ground, ifaces []iface.Interface) {
	for _, i := range ifaces {
		methods := make(map[string]rule.MethodInfo, len(i.Methods))
		for _, name := range i.Methods {
			methods[name] = rule.MethodInfo{}
		}
		g.Models[i.Name] = rule.ModelInfo{Name: i.Name, Methods: methods}
	}
}
