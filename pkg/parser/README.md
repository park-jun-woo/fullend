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
| `scenario/` | `.hurl` | `[]HurlEntry` | 통합 테스트 시나리오 파싱 |
| `manifest/` | `fullend.yaml` | `*ProjectConfig` | 프로젝트 설정 파싱 (yaml.v3 활용) |
| `toulmin/` | Go `.go` | `*Graph` | Toulmin 규칙 그래프 파싱 (Go AST 활용) |

## 외부 라이브러리 (래퍼 없이 직접 사용)

| 라이브러리 | 대상 | 출력 | 설명 |
|-----------|------|------|------|
| `kin-openapi` | OpenAPI YAML | `*openapi3.T` | API 엔드포인트 스키마 |
| `pg_query_go` | SQL DDL | `*pg_query.ParseResult` | 테이블/컬럼/제약조건 |
| `opa/ast` | OPA Rego | `*ast.Module` | 인가 정책 |
