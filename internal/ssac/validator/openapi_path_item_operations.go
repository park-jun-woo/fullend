//ff:method feature=symbol type=model
//ff:what openAPIPathItem.operations() 메서드
package validator

func (p openAPIPathItem) operations() []*openAPIOperation {
	var ops []*openAPIOperation
	for _, op := range []*openAPIOperation{p.Get, p.Post, p.Put, p.Delete} {
		if op != nil {
			ops = append(ops, op)
		}
	}
	return ops
}
