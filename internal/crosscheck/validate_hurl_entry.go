//ff:func feature=crosscheck type=rule control=sequence topic=scenario-check
//ff:what 단일 Hurl 항목을 OpenAPI 라우트와 대조 검증
package crosscheck

import "fmt"

// validateHurlEntry validates a single hurl entry against OpenAPI routes.
func validateHurlEntry(e hurlEntry, routes []apiRoute) []CrossError {
	hurlSegs := normalizeHurlPath(e.Path)

	matched, matchedRoute := findMatchingRoute(hurlSegs, e.Method, routes)

	if !matched {
		return []CrossError{{
			Rule:       "Hurl→OpenAPI",
			Level:      "ERROR",
			Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
			Message:    fmt.Sprintf("path %s not found in OpenAPI", e.Path),
			Suggestion: "OpenAPI에 해당 경로를 추가하거나, .hurl 파일의 URL을 수정하세요",
		}}
	}

	if matchedRoute == nil {
		return []CrossError{{
			Rule:       "Hurl→OpenAPI",
			Level:      "ERROR",
			Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
			Message:    fmt.Sprintf("method %s %s — path exists but method not defined in OpenAPI", e.Method, e.Path),
			Suggestion: "OpenAPI에 해당 메서드를 추가하거나, .hurl 파일의 HTTP 메서드를 수정하세요",
		}}
	}

	if e.StatusCode != "" && !matchedRoute.Responses[e.StatusCode] {
		return []CrossError{{
			Rule:       "Hurl→OpenAPI",
			Level:      "WARNING",
			Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
			Message:    fmt.Sprintf("status %s for %s %s not defined in OpenAPI responses", e.StatusCode, e.Method, e.Path),
			Suggestion: "OpenAPI responses에 해당 상태코드를 추가하거나, .hurl 파일의 기대 상태코드를 확인하세요",
		}}
	}

	return nil
}
