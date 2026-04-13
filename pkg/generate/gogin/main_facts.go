//ff:type feature=gen-gogin type=model topic=main-init
//ff:what MainFacts — DecideMainInit 입력 사실 묶음

package gogin

import (
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// MainFacts carries the precomputed facts for DecideMainInit.
type MainFacts struct {
	ServiceFuncs   []ssacparser.ServiceFunc
	Domains        []string
	QueueBackend   string
	SessionBackend string
	CacheBackend   string
	FileConfig     *manifest.FileBackend
}
