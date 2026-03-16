//ff:type feature=symbol type=model topic=ddl
//ff:what 테이블 인덱스
package validator

// Index는 테이블 인덱스다.
type Index struct {
	Name     string   // 인덱스 이름 (e.g. "idx_reservations_room_time")
	Columns  []string // 인덱스 컬럼 목록
	IsUnique bool     // UNIQUE INDEX 또는 UNIQUE 제약
}
