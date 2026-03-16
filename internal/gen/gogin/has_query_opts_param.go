//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=query-opts
//ff:what reports whether the method has a QueryOpts parameter

package gogin

// hasQueryOptsParam reports whether the method has a QueryOpts parameter.
func hasQueryOptsParam(m ifaceMethod) bool {
	for _, p := range m.Params {
		if p.Type == "QueryOpts" {
			return true
		}
	}
	return false
}
