✅ 완료

# Phase031: filefunc 도입 — 1파일 1개념 구조 전환

## 목표

fullend 코드베이스에 [filefunc](https://github.com/park-jun-woo/filefunc) 규칙을 적용하여, AI 에이전트가 `grep → read` 시 파일 단위 노이즈를 최소화한다.

설계 원칙:
1. **점진적 전환** — 전체 리팩토링이 아닌, 비대 파일 상위 10개만 1차 분해
2. **codebook 선행** — 먼저 codebook.yaml을 정의하고, 신규/분해 파일에만 `//ff:` 어노테이션 적용
3. **기존 기능 무파괴** — `go test ./...` 통과 유지가 전제, 파일 분리만 수행
4. **산출물 제외** — `artifacts/`, `pkg/` (외부 제공 패키지)는 filefunc 규칙 적용 대상에서 제외

## 동기

현재 fullend의 비대 파일 현황:

| 파일 | 줄 | 함수 수 | 문제 |
|---|---|---|---|
| `internal/ssac/validator/symbol.go` | 1,113 | 28 | DDL 파서, OpenAPI 파서, Go 인터페이스 파서, 유틸 함수가 한 파일에 밀집 |
| `internal/gen/gogin/model_impl.go` | 1,107 | 21 | 모델별 코드젠 함수가 한 파일에 밀집 |
| `internal/gen/hurl/hurl.go` | 1,048 | 35 | Hurl 생성 로직 전체가 한 파일 |
| `internal/ssac/validator/validator.go` | 937 | 24 | 검증 규칙이 한 파일에 밀집 |
| `internal/orchestrator/validate.go` | 843 | — | validate 오케스트레이션 + 출력 포매팅 혼재 |

AI 에이전트가 `symbol.go`를 읽으면 1,113줄 중 필요한 함수는 1~2개. 나머지 1,000줄은 컨텍스트 오염이다.

## 설계

### 1단계: codebook.yaml 정의

프로젝트 루트에 `codebook.yaml` 생성:

```yaml
required:
  feature:
    - orchestrator    # CLI 파싱 + 실행 순서 제어
    - crosscheck      # SSOT 간 교차 검증
    - ssac-parse      # SSaC 파서
    - ssac-validate   # SSaC 검증
    - ssac-gen        # SSaC 코드젠
    - stml-parse      # STML 파서
    - stml-validate   # STML 검증
    - stml-gen        # STML 코드젠
    - funcspec        # Func 스펙 파서
    - gen-gogin       # Go+Gin 코드젠
    - gen-hurl        # Hurl 테스트 코드젠
    - gen-react       # React TSX 코드젠
    - genmodel        # 외부 OpenAPI → Go model 생성
    - genapi          # 코드젠 공통 인터페이스
    - projectconfig   # fullend.yaml 파서
    - statemachine    # Mermaid stateDiagram 파서
    - policy          # OPA Rego 정책 파서
    - scenario        # Hurl 시나리오 파서
    - reporter        # 검증 결과 출력
    - contract        # contract hash 계산
    - symbol          # SymbolTable — DDL/OpenAPI/Go 심볼 수집
    - cli             # CLI 엔트리포인트
  type:
    - command         # CLI 명령 실행기
    - rule            # 검증 규칙
    - parser          # 파서 함수
    - generator       # 코드 생성 함수
    - model           # 데이터 구조체
    - formatter       # 출력 포매터
    - loader          # 파일/설정 로더
    - util            # 유틸리티 함수
    - walker          # 파일 순회

optional:
  ssot:
    - openapi
    - ddl
    - ssac
    - stml
    - states
    - policy
    - scenario
    - funcspec
    - config
  pattern:
    - error-collection
    - early-return
    - symbol-table
```

### 2단계: filefunc 설치 + Makefile 통합

```makefile
# Makefile 추가
.PHONY: ff-validate
ff-validate:
	filefunc validate ./internal/

.PHONY: ff-annotate
ff-annotate:
	filefunc annotate ./internal/
```

CI (`go test` 전)에 `filefunc validate` 추가는 3단계 완료 후.

### 3단계: 비대 파일 분해 — 1차 (상위 5개)

#### 3-1. `internal/ssac/validator/symbol.go` (1,113줄 → 10+개 파일)

현재 28개 함수를 기능별로 분리:

| 신규 파일 | 원본 함수 | feature | type |
|---|---|---|---|
| `symbol_table.go` | `SymbolTable` struct, `Clone()` | symbol | model |
| `load_symbol_table.go` | `LoadSymbolTable()` | symbol | loader |
| `load_openapi.go` | `loadOpenAPI()` | symbol | loader |
| `load_ddl.go` | `loadDDL()` | symbol | loader |
| `parse_ddl_tables.go` | `parseDDLTables()` | symbol | parser |
| `parse_inline_fk.go` | `parseInlineFK()` | symbol | parser |
| `parse_constraint_fk.go` | `parseConstraintFK()` | symbol | parser |
| `parse_create_index.go` | `parseCreateIndex()` | symbol | parser |
| `load_go_interfaces.go` | `loadGoInterfaces()` | symbol | loader |
| `load_sqlc_queries.go` | `loadSqlcQueries()` | symbol | loader |
| `collect_schema_fields.go` | `collectSchemaFields()` | symbol | util |
| `pg_type_to_go.go` | `pgTypeToGo()` | symbol | util |
| `oa_type_to_go.go` | `oaTypeToGo()` | symbol | util |
| `openapi_types.go` | `openAPISpec`, `openAPISchema` 등 내부 YAML struct | symbol | model |
| `model_symbol.go` | `ModelSymbol`, `MethodInfo`, `OperationSymbol` 등 | symbol | model |
| `ddl_table.go` | `DDLTable`, `ForeignKey`, `Index`, `PathParam` 등 | symbol | model |

**예외 처리**: `openapi_types.go`에는 `openAPISpec`, `openAPISchema`, `openAPIPathItem`, `openAPIOperation`, `openAPIParameter`, `openAPIRequestBody`, `openAPIResponse`, `openAPIMediaType` 등 관련 struct가 밀집 — 이들은 의미적으로 한 묶음이므로 F2 예외(F7: semantically grouped) 적용, 1파일에 유지.

#### 3-2. `internal/gen/gogin/model_impl.go` (1,107줄 → 6+개 파일)

모델별 코드젠 함수를 분리:

| 신규 파일 | 원본 함수 | feature | type |
|---|---|---|---|
| `gen_model_impl.go` | `generateModelImpl()` 메인 디스패처 | gen-gogin | generator |
| `gen_model_scan.go` | `generateScanFunc()` | gen-gogin | generator |
| `gen_model_create.go` | `generateCreateMethod()` | gen-gogin | generator |
| `gen_model_find.go` | `generateFindMethod()` | gen-gogin | generator |
| `gen_model_list.go` | `generateListMethod()` | gen-gogin | generator |
| `gen_model_update.go` | `generateUpdateMethod()` | gen-gogin | generator |
| `gen_model_with_tx.go` | `generateWithTxMethod()` | gen-gogin | generator |

#### 3-3. `internal/gen/hurl/hurl.go` (1,048줄 → 8+개 파일)

| 신규 파일 | 원본 함수 | feature | type |
|---|---|---|---|
| `generate_hurl_tests.go` | `generateHurlTests()` 메인 | gen-hurl | generator |
| `gen_hurl_auth.go` | 인증 시나리오 생성 | gen-hurl | generator |
| `gen_hurl_crud.go` | CRUD 시나리오 생성 | gen-hurl | generator |
| `gen_hurl_state.go` | 상태 전이 시나리오 생성 | gen-hurl | generator |
| `gen_hurl_invariant.go` | invariant 시나리오 생성 | gen-hurl | generator |
| `hurl_assertion.go` | assertion 빌더 | gen-hurl | util |
| `hurl_request.go` | request 빌더 | gen-hurl | util |
| `hurl_capture.go` | capture 빌더 | gen-hurl | util |

#### 3-4. `internal/ssac/validator/validator.go` (937줄 → 7+개 파일)

검증 규칙별로 분리:

| 신규 파일 | 원본 함수 | feature | type |
|---|---|---|---|
| `validate.go` | `Validate()` 메인 디스패처 | ssac-validate | command |
| `check_model_ref.go` | 모델 참조 검증 | ssac-validate | rule |
| `check_method_ref.go` | 메서드 참조 검증 | ssac-validate | rule |
| `check_response.go` | @response 검증 | ssac-validate | rule |
| `check_auth.go` | @auth 검증 | ssac-validate | rule |
| `check_state.go` | @state 검증 | ssac-validate | rule |
| `check_call.go` | @call 검증 | ssac-validate | rule |

#### 3-5. `internal/orchestrator/validate.go` (843줄)

| 신규 파일 | 원본 함수 | feature | type |
|---|---|---|---|
| `validate.go` | `ValidateWith()` 메인 오케스트레이터 | orchestrator | command |
| `run_cross_validate.go` | `runCrossValidate()` | orchestrator | command |
| `format_errors.go` | 에러 출력 포매팅 | orchestrator | formatter |

### 4단계: 분해 파일에 `//ff:` 어노테이션 추가

분해 시 각 파일 상단에 어노테이션 추가. 패턴:

```go
//ff:func feature=symbol type=parser
//ff:what DDL CREATE TABLE 문에서 컬럼명, 타입, FK, 인덱스를 추출
package validator

func parseDDLTables(content string, tables map[string]DDLTable) {
```

```go
//ff:type feature=symbol type=model
//ff:what DDL 테이블의 컬럼·FK·인덱스·PK 정보를 담는 구조체
package validator

type DDLTable struct {
```

### 5단계: `filefunc annotate` 실행

분해 완료 후 `filefunc annotate ./internal/`로 `//ff:calls`, `//ff:uses` 자동 생성.

### 6단계: `filefunc validate` CI 통합

`.github/workflows/` 또는 `Makefile`의 `test` 타겟에 추가:

```makefile
.PHONY: test
test:
	filefunc validate ./internal/
	go test ./...
```

**신규 파일에만 강제**: 기존 미분해 파일은 2차 분해 전까지 `//ff:` 어노테이션 없이 유지. `filefunc validate`는 어노테이션이 없는 파일을 A1 ERROR로 보고하므로, 1차 분해 대상 외 파일은 `.ffignore` 또는 디렉토리 단위 제외로 처리.

**제외 대상**:
- `artifacts/` — 코드젠 산출물
- `pkg/` — 외부 제공 패키지 (filefunc 규칙과 무관)
- `cmd/` — main.go 1개뿐, 분해 불필요
- `*_test.go` — F5 예외로 복수 함수 허용

## 변경 파일

| 파일 | 변경 |
|---|---|
| `codebook.yaml` | **신규** — fullend 프로젝트 codebook 정의 |
| `Makefile` | `ff-validate`, `ff-annotate` 타겟 추가 |
| `internal/ssac/validator/symbol.go` | 28개 함수 → 16개 파일로 분해, 원본 삭제 |
| `internal/ssac/validator/symbol_table.go` | **신규** — SymbolTable struct + Clone() |
| `internal/ssac/validator/load_openapi.go` | **신규** — loadOpenAPI() |
| `internal/ssac/validator/load_ddl.go` | **신규** — loadDDL() |
| `internal/ssac/validator/parse_ddl_tables.go` | **신규** — parseDDLTables() |
| `internal/ssac/validator/openapi_types.go` | **신규** — openAPI 내부 YAML struct 묶음 |
| `internal/ssac/validator/model_symbol.go` | **신규** — ModelSymbol, MethodInfo 등 |
| `internal/ssac/validator/ddl_table.go` | **신규** — DDLTable, ForeignKey 등 |
| `internal/ssac/validator/` (기타) | **신규** — load_*, parse_*, *_to_go.go 등 |
| `internal/gen/gogin/model_impl.go` | 21개 함수 → 7개 파일로 분해, 원본 삭제 |
| `internal/gen/gogin/gen_model_*.go` | **신규** — 모델별 코드젠 분리 파일 |
| `internal/gen/hurl/hurl.go` | 35개 함수 → 8개 파일로 분해, 원본 삭제 |
| `internal/gen/hurl/gen_hurl_*.go` | **신규** — 시나리오별 코드젠 분리 파일 |
| `internal/ssac/validator/validator.go` | 24개 함수 → 7개 파일로 분해, 원본 삭제 |
| `internal/ssac/validator/check_*.go` | **신규** — 검증 규칙별 분리 파일 |
| `internal/orchestrator/validate.go` | 분해 → validate.go + run_cross_validate.go + format_errors.go |

**예상 파일 수 변화**: ~80개 `.go` → ~130개 `.go` (1차 분해 5개 파일 대상, +50개)

## 의존성

- Phase030(커스텀 JWT Claims) 이후. Phase032(OpenAPI Validation Tags)와 독립이나, Phase032가 `symbol.go`에 함수를 추가하므로 Phase031을 먼저 완료하면 Phase032의 신규 함수가 자연스럽게 별도 파일로 생성됨.
- 외부 도구: `filefunc` CLI (`go install github.com/park-jun-woo/filefunc/cmd/filefunc@latest`)

## 검증

1. `go test ./...` — 분해 후 모든 기존 테스트 통과 (기능 무파괴 확인)
2. `go build ./cmd/fullend/` — 빌드 성공
3. `go vet ./...` — 경고 없음
4. `filefunc validate ./internal/ssac/validator/` — 분해된 파일이 F1, F2, A1, A3 통과
5. `filefunc validate ./internal/gen/gogin/` — 분해된 파일 통과
6. `filefunc validate ./internal/gen/hurl/` — 분해된 파일 통과
7. `filefunc annotate ./internal/` → `//ff:calls`, `//ff:uses` 자동 생성 확인
8. `fullend validate specs/dummys/zenflow-try05/` — 기존 기능 정상 동작
9. `fullend gen specs/dummys/zenflow-try05/ artifacts/dummys/zenflow-try05/` — 코드젠 정상

## 리스크

- **파일 수 폭발** — 80 → 130개는 관리 가능 범위. 단, 2차 분해(나머지 파일)까지 하면 200+개 예상. 이는 filefunc의 설계 의도(파일 수 증가는 feature)와 일치.
- **import cycle** — 같은 패키지 내 분해이므로 import cycle 위험 없음. 패키지 경계를 넘는 분해는 하지 않음.
- **git diff 폭발** — 파일 이동이 많아 PR diff가 거대해짐. 파일별로 `git mv` + 수정을 분리하여 리뷰 용이성 확보. 또는 패키지 단위로 분해 커밋을 나눔 (5개 커밋).
- **테스트 파일** — `validator_test.go` (1,384줄, 84 함수)는 F5 예외로 분해 대상 아님. 향후 테스트 가독성 개선 시 별도 분해 가능.
- **codebook 진화** — feature/type 값은 프로젝트 성장에 따라 추가됨. codebook 변경 시 기존 어노테이션과 불일치 가능 → `filefunc validate`가 즉시 감지.
- **2차 분해 시점** — 1차 분해 후 나머지 중간 크기 파일(400~700줄)은 Phase031에서 다루지 않음. 필요 시 별도 Phase로.
