//ff:type feature=gen-gogin type=model topic=main-init
//ff:what BackendNeed — session/cache/file 단일 백엔드 활성화 상태

package gogin

import "github.com/park-jun-woo/fullend/pkg/parser/manifest"

// BackendNeed expresses whether a backed feature is enabled and which backend.
type BackendNeed struct {
	Enabled    bool
	Backend    string                // "postgres"/"memory" for session/cache; "local"/"s3" for file
	FileConfig *manifest.FileBackend // non-nil only for File
}
