//ff:type feature=crosscheck type=model topic=ddl-coverage
//ff:what DDL @archived 태그 정보를 담는 구조체
package crosscheck

// ArchivedInfo holds @archived tags parsed from DDL files.
type ArchivedInfo struct {
	Tables  map[string]bool            // "legacy_notifications" → true
	Columns map[string]map[string]bool // "courses" → {"old_category": true}
}
