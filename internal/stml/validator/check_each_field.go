//ff:func feature=stml-validate type=rule control=sequence
//ff:what 단일 data-each 필드가 응답에 존재하고 배열인지 확인
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func checkEachField(e parser.EachBlock, opID, file string, api APISymbol) *ValidationError {
	fs, ok := api.ResponseFields[e.Field]
	if !ok {
		err := errEachNotFound(file, opID, e.Field)
		return &err
	}
	if fs.Type != "array" {
		err := errEachNotArray(file, opID, e.Field)
		return &err
	}
	return nil
}
