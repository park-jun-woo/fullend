//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what converts OpenAPI path params {Name} to gin style :Name

package gogin

import "strings"

// convertPathParamsGin converts OpenAPI path params {Name} to gin style :Name.
func convertPathParamsGin(path string) string {
	result := path
	for {
		start := strings.Index(result, "{")
		if start < 0 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end < 0 {
			break
		}
		paramName := result[start+1 : start+end]
		result = result[:start] + ":" + paramName + result[start+end+1:]
	}
	return result
}
