//ff:func feature=gen-gogin type=util control=sequence topic=ddl
//ff:what parseSeedKey — "<ref>#<N>" 키를 (refTable, id) 로 분해

package db

import (
	"strconv"
	"strings"
)

func parseSeedKey(key string) (refTable string, id int64) {
	i := strings.LastIndex(key, "#")
	if i < 0 {
		return key, 0
	}
	n, _ := strconv.ParseInt(key[i+1:], 10, 64)
	return key[:i], n
}
