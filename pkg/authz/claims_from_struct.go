//ff:func feature=pkg-authz type=util control=iteration dimension=1
//ff:what ClaimsFromStruct — struct/pointer 를 snake_case key 의 map[string]any 로 변환

package authz

import "reflect"

// ClaimsFromStruct converts a struct (or pointer to struct) into a map used as
// OPA `input.claims`. Key resolution order per field:
//  1. `authz:"<key>"` struct tag (JWT claim key from manifest.auth.claims config)
//  2. fallback: snake_case version of field name
//
// "ID"(tag:user_id) → "user_id"; "OrgID"(no tag) → "org_id"; "Role" → "role".
func ClaimsFromStruct(s any) map[string]any {
	if s == nil {
		return nil
	}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	out := make(map[string]any, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() {
			continue
		}
		key := fld.Tag.Get("authz")
		if key == "" {
			key = toSnakeCase(fld.Name)
		}
		out[key] = v.Field(i).Interface()
	}
	return out
}

