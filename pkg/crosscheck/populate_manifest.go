//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateManifest — fullend.yaml에서 middleware, claims, roles, queue 설정 추출
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateManifest(g *rule.Ground, fs *fullend.Fullstack) {
	if fs.Manifest == nil {
		return
	}
	middleware := make(rule.StringSet)
	for _, m := range fs.Manifest.Backend.Middleware {
		middleware[m] = true
	}
	g.Lookup["Config.middleware"] = middleware

	if fs.Manifest.Backend.Auth != nil {
		claims := make(rule.StringSet)
		claimKeys := make(rule.StringSet)
		for field, def := range fs.Manifest.Backend.Auth.Claims {
			claims[field] = true
			claimKeys[def.Key] = true
		}
		g.Lookup["Config.claims"] = claims
		g.Lookup["Config.claims.keys"] = claimKeys
		g.Config["backend.auth.claims"] = len(claims) > 0

		roles := make(rule.StringSet)
		for _, r := range fs.Manifest.Backend.Auth.Roles {
			roles[r] = true
		}
		g.Lookup["Config.roles"] = roles
	}

	if fs.Manifest.Queue != nil && fs.Manifest.Queue.Backend != "" {
		g.Config["queue.backend"] = true
	}
	if len(fs.Manifest.Backend.Middleware) > 0 {
		g.Config["backend.middleware"] = true
	}
}
