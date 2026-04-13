//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what seedAlreadyExists — 테이블 Seeds 에 id=N 행이 존재하는지

package db

import (
	"strconv"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

func seedAlreadyExists(t *ddl.Table, id int64) bool {
	idStr := strconv.FormatInt(id, 10)
	for _, row := range t.Seeds {
		if row["id"] == idStr {
			return true
		}
	}
	return false
}
