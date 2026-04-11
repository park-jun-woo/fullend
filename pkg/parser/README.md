# pkg/parser

## 통합

| 패키지 | 설명 |
|--------|------|
| `fullend/` | 모든 SSOT 파싱 결과를 담는 `Fullstack` 컨테이너 + `ParseAll()` |

## 자체 파서

| 패키지 | 입력 | 출력 | 설명 |
|--------|------|------|------|
| `ssac/` | `.ssac` | `[]ServiceFunc` | 서비스 함수 시퀀스 파싱 (Go AST 활용) |
| `stml/` | `.html` | `[]PageSpec` | 프론트엔드 페이지 템플릿 파싱 (x/net/html 활용) |
| `statemachine/` | Mermaid `.md` | `[]*StateDiagram` | 상태 전이 다이어그램 파싱 |
| `funcspec/` | Go `.go` | `[]FuncSpec` | 커스텀 함수 스펙 파싱 (Go AST 활용) |
| `hurl/` | `.hurl` | `[]HurlEntry` | 통합 테스트 시나리오 파싱 |
| `manifest/` | `fullend.yaml` | `*ProjectConfig` | 프로젝트 설정 파싱 (yaml.v3 활용) |
| `toulmin/` | Go `.go` | `*Graph` | Toulmin 규칙 그래프 파싱 (Go AST 활용) |

## 구조화 파서 + 외부 검증

| 패키지 | 구조화 파서 | 출력 | 외부 검증 | 설명 |
|--------|-----------|------|----------|------|
| `ddl/` | `ParseTables(dir)` | `[]Table` | `ParseDir()` → `pg_query_go` | DDL 테이블/컬럼/FK/인덱스/CHECK |
| `rego/` | `ParsePolicies(dir)` | `[]Policy` | `ParseDir()` → `opa/ast` | allow 규칙, @ownership, claims 참조 |

## 외부 라이브러리 (래퍼 없이 직접 사용)

| 라이브러리 | 대상 | 출력 | 설명 |
|-----------|------|------|------|
| `kin-openapi` | OpenAPI YAML | `*openapi3.T` | API 엔드포인트 스키마 |

## DDL 구조체

```go
type Table struct {
    Name        string
    Columns     map[string]string   // column → Go type
    ColumnOrder []string
    ForeignKeys []ForeignKey        // {Column, RefTable, RefColumn}
    Indexes     []Index             // {Name, Columns, IsUnique}
    PrimaryKey  []string
    VarcharLen  map[string]int      // column → VARCHAR(N)
    CheckEnums  map[string][]string // column → CHECK IN values
}
```

## Rego 구조체

```go
type Policy struct {
    File       string
    Rules      []AllowRule          // {Actions, Resource, UsesOwner, UsesRole, RoleValue}
    Ownerships []OwnershipMapping   // {Resource, Table, Column, JoinTable, JoinFK}
    ClaimsRefs []string             // input.claims.xxx 참조 (중복 제거)
}
```
