//ff:func feature=crosscheck type=util control=selection
//ff:what kin-openapi 확장 값을 JSON 언마샬
package crosscheck

import "encoding/json"

// unmarshalExt handles kin-openapi extension values which may be json.RawMessage.
func unmarshalExt(v any, dst any) error {
	switch val := v.(type) {
	case json.RawMessage:
		return json.Unmarshal(val, dst)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, dst)
	}
}
