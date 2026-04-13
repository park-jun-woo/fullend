//ff:func feature=gen-gogin type=decider control=sequence topic=main-init
//ff:what DecideMainInit — MainFacts → InitNeeds (6축 독립 판정, depth 1)

package gogin

// DecideMainInit evaluates each axis independently (depth 1 per axis).
func DecideMainInit(f MainFacts) InitNeeds {
	needs := InitNeeds{}

	if anyDomainNeedsAuth(f.ServiceFuncs, f.Domains) {
		needs.Auth = true
		needs.Authz = true
	}

	if f.QueueBackend != "" && (len(collectSubscribers(f.ServiceFuncs)) > 0 || hasPublishSequence(f.ServiceFuncs)) {
		needs.Queue = true
	}

	if f.SessionBackend == "postgres" || f.SessionBackend == "memory" {
		needs.Session = BackendNeed{Enabled: true, Backend: f.SessionBackend}
	}
	if f.CacheBackend == "postgres" || f.CacheBackend == "memory" {
		needs.Cache = BackendNeed{Enabled: true, Backend: f.CacheBackend}
	}
	if f.FileConfig != nil {
		needs.File = BackendNeed{
			Enabled:    true,
			Backend:    f.FileConfig.Backend,
			FileConfig: f.FileConfig,
		}
	}

	needs.NeedsContextImport =
		(needs.Session.Enabled && needs.Session.Backend == "postgres") ||
		(needs.Cache.Enabled && needs.Cache.Backend == "postgres") ||
		(needs.File.Enabled && needs.File.Backend == "s3") ||
		needs.Queue

	return needs
}
