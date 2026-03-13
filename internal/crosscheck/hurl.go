package crosscheck

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	reHurlRequest  = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+\{\{host\}\}(.+)`)
	reHurlResponse = regexp.MustCompile(`^HTTP\s+(\d+)`)
)

// hurlEntry represents one request/response pair extracted from a .hurl file.
type hurlEntry struct {
	Method     string
	Path       string
	StatusCode string
	File       string
	Line       int
}

// CheckHurlFiles validates that .hurl scenario files reference valid OpenAPI paths.
func CheckHurlFiles(hurlFiles []string, doc *openapi3.T) []CrossError {
	var errs []CrossError

	// Build normalized OpenAPI path set: segment arrays.
	type apiRoute struct {
		Method    string
		Segments  []string
		Responses map[string]bool // status codes defined
	}
	var routes []apiRoute
	if doc.Paths != nil {
		for path, pi := range doc.Paths.Map() {
			segs := normalizeOpenAPIPath(path)
			for method, op := range pi.Operations() {
				responseCodes := make(map[string]bool)
				if op.Responses != nil {
					for code := range op.Responses.Map() {
						responseCodes[code] = true
					}
				}
				routes = append(routes, apiRoute{
					Method:    strings.ToUpper(method),
					Segments:  segs,
					Responses: responseCodes,
				})
			}
		}
	}

	// Parse .hurl files and check each entry.
	for _, f := range hurlFiles {
		entries := parseHurlFile(f)
		for _, e := range entries {
			hurlSegs := normalizeHurlPath(e.Path)

			// Find matching route.
			matched := false
			var matchedRoute *apiRoute
			for i := range routes {
				if segmentsMatch(hurlSegs, routes[i].Segments) {
					matched = true
					if routes[i].Method == e.Method {
						matchedRoute = &routes[i]
						break
					}
				}
			}

			if !matched {
				errs = append(errs, CrossError{
					Rule:       "Hurl→OpenAPI",
					Level:      "ERROR",
					Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
					Message:    fmt.Sprintf("path %s not found in OpenAPI", e.Path),
					Suggestion: "OpenAPI에 해당 경로를 추가하거나, .hurl 파일의 URL을 수정하세요",
				})
				continue
			}

			if matchedRoute == nil {
				errs = append(errs, CrossError{
					Rule:       "Hurl→OpenAPI",
					Level:      "ERROR",
					Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
					Message:    fmt.Sprintf("method %s %s — path exists but method not defined in OpenAPI", e.Method, e.Path),
					Suggestion: "OpenAPI에 해당 메서드를 추가하거나, .hurl 파일의 HTTP 메서드를 수정하세요",
				})
				continue
			}

			// Check status code (WARNING level).
			if e.StatusCode != "" && !matchedRoute.Responses[e.StatusCode] {
				errs = append(errs, CrossError{
					Rule:       "Hurl→OpenAPI",
					Level:      "WARNING",
					Context:    fmt.Sprintf("%s:%d", e.File, e.Line),
					Message:    fmt.Sprintf("status %s for %s %s not defined in OpenAPI responses", e.StatusCode, e.Method, e.Path),
					Suggestion: "OpenAPI responses에 해당 상태코드를 추가하거나, .hurl 파일의 기대 상태코드를 확인하세요",
				})
			}
		}
	}

	return errs
}

// parseHurlFile extracts request/response pairs from a .hurl file.
func parseHurlFile(path string) []hurlEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var entries []hurlEntry
	var current *hurlEntry

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if m := reHurlRequest.FindStringSubmatch(line); m != nil {
			// Flush previous entry.
			if current != nil {
				entries = append(entries, *current)
			}
			current = &hurlEntry{
				Method: m[1],
				Path:   strings.TrimSpace(m[2]),
				File:   path,
				Line:   lineNum,
			}
			continue
		}

		if m := reHurlResponse.FindStringSubmatch(line); m != nil && current != nil {
			current.StatusCode = m[1]
		}
	}

	if current != nil {
		entries = append(entries, *current)
	}

	return entries
}

// normalizeHurlPath converts a Hurl URL path to normalized segments.
// {{variable}} → ":param", pure numeric literals (e.g. "999999") → ":param"
func normalizeHurlPath(path string) []string {
	path = strings.TrimSpace(path)
	// Remove query string.
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var segs []string
	reVar := regexp.MustCompile(`\{\{.+?\}\}`)
	reNumeric := regexp.MustCompile(`^\d+$`)
	for _, p := range parts {
		if p == "" {
			continue
		}
		if reVar.MatchString(p) || reNumeric.MatchString(p) {
			segs = append(segs, ":param")
		} else {
			segs = append(segs, p)
		}
	}
	return segs
}

// normalizeOpenAPIPath converts an OpenAPI path to normalized segments.
// {param} → ":param"
func normalizeOpenAPIPath(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var segs []string
	reParam := regexp.MustCompile(`^\{.+\}$`)
	for _, p := range parts {
		if p == "" {
			continue
		}
		if reParam.MatchString(p) {
			segs = append(segs, ":param")
		} else {
			segs = append(segs, p)
		}
	}
	return segs
}

// segmentsMatch checks if two segment arrays match.
func segmentsMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
