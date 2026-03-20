# ✅ Phase001 — 파서 에러 라인 정보 통합

## 목표

모든 파서가 에러 발생 시 `pkg/diagnostic.Diagnostic`으로 파일 경로 + 라인 번호를 반환하도록 통일한다.

## 현황

| 파서 | 라인 정보 | 상태 |
|------|-----------|------|
| `ssac/` | Go AST `token.Pos`로 추적 가능 | 부분 구현 |
| `stml/` | `x/net/html`이 위치 제공 안 함 | 미구현 |
| `statemachine/` | 정규식 라인 파싱, 라인 번호 미추적 | 미구현 |
| `funcspec/` | Go AST 활용하지만 에러에 라인 미포함 | 미구현 |
| `scenario/` | `HurlEntry.Line`은 데이터 위치, 에러 Diagnostic 미구현 | 미구현 |
| `manifest/` | `yaml.v3` 에러에 라인 포함 | 부분 구현 |
| `ddl/` | `pg_query_go` 에러에 Cursorpos 포함 | 미구현 |
| `rego/` | `opa/ast` 에러에 Location{Row,Col} 포함 | 미구현 |

## 변경 파일 목록

### 공통
- `pkg/diagnostic/diagnostic.go` — 이미 생성됨, 변경 없음

### ssac 파서
- `pkg/parser/ssac/parse_file.go` — 에러 메시지에 라인 번호 통일
- `pkg/parser/ssac/parse_func_decl.go` — `fset.Position()` 활용 강화

### stml 파서
- `pkg/parser/stml/parse_reader.go` — 라인 카운터 도입 (토크나이저 래핑)
- `pkg/parser/stml/type_page_spec.go` — 에러 반환 타입에 Diagnostic 추가 검토
- `pkg/parser/stml/parse_fetch_block.go` — 에러 시 라인 정보 포함
- `pkg/parser/stml/parse_action_block.go` — 에러 시 라인 정보 포함

### statemachine 파서
- `pkg/parser/statemachine/parse.go` — 라인 순회 시 lineNum 추적, 에러에 포함

### funcspec 파서
- `pkg/parser/funcspec/parse_file.go` — `fset.Position()`으로 라인 번호 추출
- `pkg/parser/funcspec/parse_comment_group.go` — 어노테이션 에러에 라인 포함

### scenario 파서
- `pkg/parser/scenario/parse_hurl_file.go` — 구문 에러 시 Diagnostic 반환 추가

### manifest 파서
- `pkg/parser/manifest/load.go` — yaml.v3 에러에서 라인 추출하여 Diagnostic 변환

### ddl 래퍼
- `pkg/parser/ddl/parse_dir.go` — `pg_query.Error`의 Cursorpos → Diagnostic 변환은 ddl 패키지 책임

### rego 래퍼
- `pkg/parser/rego/parse_dir.go` — `ast.Error`의 Location → Diagnostic 변환은 rego 패키지 책임

## 의존성

- `pkg/diagnostic` (이미 생성됨)

## 구현 방침

1. `pkg/diagnostic/`는 타입 정의만. 변환 로직 없음
2. 각 파서가 자기 에러를 `Diagnostic`으로 변환하는 책임을 진다
   - `ddl/` → `pg_query.Error` → `Diagnostic`
   - `rego/` → `ast.Error` → `Diagnostic`
   - `manifest/` → `yaml.TypeError` → `Diagnostic`
   - `ssac/`, `funcspec/` → Go AST `fset.Position()` → `Diagnostic`
   - `stml/` → `\n` 카운트 근사치 → `Diagnostic`
   - `statemachine/` → 라인 카운터 → `Diagnostic`
   - `scenario/` → 라인 카운터 → `Diagnostic`
3. 파서의 반환 타입은 기존 유지 ([]ServiceFunc 등). 에러만 `[]diagnostic.Diagnostic`으로 통일
4. `HurlEntry.Line` 같은 데이터 위치 정보와 파싱 에러는 별개. 데이터 위치는 AST에, 에러는 Diagnostic에
5. 구문 에러(syntax error)만 파서가 담당. 의미 에러는 validate 단계
6. 라인 번호를 모를 경우 `Line: 0`으로 표기 (unknown)
7. `x/net/html`은 라인 정보가 없으므로, 토크나이저를 래핑하여 `\n` 카운트로 근사치 추적

## 검증 방법

- 각 파서별 테스트에서 의도적 구문 에러 입력 → Diagnostic의 File, Line, Message 검증
- `go test ./pkg/parser/...` 전체 통과
- `filefunc validate` 위반 0건
