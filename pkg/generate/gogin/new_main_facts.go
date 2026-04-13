//ff:func feature=gen-gogin type=util control=sequence topic=main-init
//ff:what newMainFacts — MainGenInput → MainFacts 프로젝션

package gogin

// newMainFacts projects MainGenInput into MainFacts.
func newMainFacts(in MainGenInput) MainFacts {
	return MainFacts{
		ServiceFuncs:   in.ServiceFuncs,
		Domains:        uniqueDomains(in.ServiceFuncs),
		QueueBackend:   in.QueueBackend,
		SessionBackend: in.SessionBackend,
		CacheBackend:   in.CacheBackend,
		FileConfig:     in.FileConfig,
	}
}
