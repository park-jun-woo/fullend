//ff:func feature=rule type=util control=sequence
//ff:what extractPagination/Sort/Filter — OpenAPI x-확장에서 구조 추출
package ground

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func extractPagination(op *openapi3.Operation) *rule.PaginationSpec {
	raw, ok := op.Extensions["x-pagination"]
	if !ok {
		return nil
	}
	var parsed struct {
		Style        string `json:"style"`
		DefaultLimit int    `json:"defaultLimit"`
		MaxLimit     int    `json:"maxLimit"`
	}
	if !unmarshalExtension(raw, &parsed) {
		return nil
	}
	return &rule.PaginationSpec{
		Style:        parsed.Style,
		DefaultLimit: parsed.DefaultLimit,
		MaxLimit:     parsed.MaxLimit,
	}
}

func extractSort(op *openapi3.Operation) *rule.SortSpec {
	raw, ok := op.Extensions["x-sort"]
	if !ok {
		return nil
	}
	var parsed struct {
		Allowed   []string `json:"allowed"`
		Default   string   `json:"default"`
		Direction string   `json:"direction"`
	}
	if !unmarshalExtension(raw, &parsed) {
		return nil
	}
	return &rule.SortSpec{
		Allowed:   parsed.Allowed,
		Default:   parsed.Default,
		Direction: parsed.Direction,
	}
}

func extractFilter(op *openapi3.Operation) *rule.FilterSpec {
	raw, ok := op.Extensions["x-filter"]
	if !ok {
		return nil
	}
	var parsed struct {
		Allowed []string `json:"allowed"`
	}
	if !unmarshalExtension(raw, &parsed) {
		return nil
	}
	return &rule.FilterSpec{Allowed: parsed.Allowed}
}

func unmarshalExtension(raw any, out any) bool {
	data, err := json.Marshal(raw)
	if err != nil {
		return false
	}
	return json.Unmarshal(data, out) == nil
}
