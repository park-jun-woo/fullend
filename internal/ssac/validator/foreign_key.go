//ff:type feature=symbol type=model topic=ddl
//ff:what 외래 키 관계
package validator

// ForeignKey는 외래 키 관계다.
type ForeignKey struct {
	Column    string // 이 테이블의 컬럼 (e.g. "user_id")
	RefTable  string // 참조 테이블 (e.g. "users")
	RefColumn string // 참조 컬럼 (e.g. "id")
}
