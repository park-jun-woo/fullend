# Phase051 — package-prefix @model 폐기, session/cache/file을 @call Func으로 전환 ✅ 완료

## 목표
SSaC의 package-prefix @model 문법(`session.Session.Set`, `cache.Cache.Get` 등)을 폐기한다.
session/cache/file은 `@call` Func으로 다루며, fullend `pkg/`에서 기본 제공하되 사용자가 프로젝트 `func/`에서 커스텀할 수 있도록 한다.

## 배경
zenflow-try06 add05에서 `@put session.Session.Set(...)` 사용 시 빌드 실패. 조사 결과 SSaC의 `@get`/`@put`/`@delete`는 DDL 모델 전용으로 설계되어 있어, DDL이 아닌 인프라 서비스(session/cache/file)를 같은 경로에 넣으면 WithTx, ctx 주입, 인터페이스 생성, 인자 순서, 타입 변환이 전부 깨진다.

**근본 원인**: 코드젠 버그가 아니라 SSaC 설계 결함. `@get`/`@put`은 DDL 모델 전용이며, session/cache/file은 인프라 서비스로 별도 경로(`@call`)가 맞다.

## 설계 방침

### @model은 DDL 전담
- `@get`/`@post`/`@put`/`@delete`는 DDL 테이블 모델만 사용
- package-prefix 문법(`pkg.Model.Method`)은 ERROR로 차단

### session/cache/file은 @call Func
- `@call session.Set({Key: ..., Value: ..., TTL: ...})`
- `@call cache.Get({Key: ...})`
- `@call file.Upload({Key: ..., Body: ...})`
- fullend `pkg/`에서 `@func` 어노테이션 + Request/Response struct 기본 제공
- Fallback chain: `specs/<project>/func/<pkg>/` → `pkg/<pkg>/` — 사용자 커스텀 우선

### Purity 규칙 재정의
모든 `@call` Func에 대해 다음 규칙을 적용한다:
- **허용**: file I/O(`io`, `bufio`, `os`), session/cache 읽기·쓰기 — 일반적인 비즈니스 로직에서 필요
- **절대 금지**: DB 직접 접근(`database/sql`, `github.com/lib/pq`, `github.com/jackc/pgx`)
- **절대 금지**: API 호출(`net/http`, `net/rpc`, `google.golang.org/grpc`)

면제 로직 불필요 — 금지 목록 자체를 DB·네트워크로 한정한다.

---

## 변경 내용

### 1. validator: package-prefix @model ERROR 차단

#### `internal/ssac/validator/validate_model.go`
`seq.Package != ""`이면 WARNING 대신 ERROR 반환.
```go
if seq.Package != "" {
    errs = append(errs, ctx.err("@"+seq.Type,
        fmt.Sprintf("%s.%s — package-prefix @model은 지원하지 않습니다. @call Func을 사용하세요", seq.Package, modelName)))
    continue
}
```

### 2. validator: package-prefix 관련 코드 삭제

| 파일 | 조치 |
|---|---|
| `internal/ssac/validator/validate_package_model.go` | 삭제 |
| `internal/ssac/validator/load_package_interfaces.go` | 삭제 |
| `internal/ssac/validator/parse_package_interfaces.go` | 삭제 |
| `internal/ssac/validator/load_package_go_interfaces.go` | 삭제 |

`validate_model.go`에서 `validatePackageModel` 호출 제거, `LoadPackageInterfaces` 호출부 제거.

#### 호출부 확인
```
filefunc chain func LoadPackageInterfaces --chon 1
filefunc chain func validatePackageModel --chon 1
```

### 3. generator: package-prefix 잔여 코드 정리

| 파일 | 조치 |
|---|---|
| `internal/ssac/generator/collect_model_usages_from_func.go` | `seq.Package != ""` 조건 제거 (validate에서 이미 ERROR) |
| `internal/ssac/generator/default_message.go` | `seq.Package` 참조 제거 |

### 4. parser: package-prefix 파싱 유지
`split_package_prefix.go`와 `Sequence.Package` 필드는 유지한다. 파서가 구문을 인식해야 validator가 의미 있는 ERROR 메시지를 출력할 수 있다.

### 5. 테스트 수정

| 파일 | 조치 |
|---|---|
| `internal/ssac/generator/go_interface_test.go` | `TestGeneratePackageModelSkipInterface` 삭제 |
| `internal/ssac/generator/go_handler_test.go` | `TestGeneratePackageModelCall` 삭제 |
| `internal/ssac/validator/validator_test.go` | package-prefix 관련 테스트를 ERROR 기대로 변경 |
| `internal/ssac/parser/parser_test.go` | 파싱 테스트 유지 (파서는 변경 없음) |

### 6. pkg/ Func 스펙 추가 — session, cache, file

`pkg/session`, `pkg/cache`, `pkg/file`에 `@func` 어노테이션 + Request/Response struct 추가.

#### `pkg/session/set.go`
```go
// @func set
// @description 세션에 key-value를 저장한다

type SetRequest struct {
    Key   string
    Value string
    TTL   int64
}

type SetResponse struct{}

func Set(req SetRequest) (SetResponse, error) {
    // 실제 구현: SessionModel.Set 호출
}
```

#### `pkg/session/get.go`
```go
// @func get
type GetRequest struct { Key string }
type GetResponse struct { Value string }
func Get(req GetRequest) (GetResponse, error) { ... }
```

#### `pkg/session/delete.go`
```go
// @func delete
type DeleteRequest struct { Key string }
type DeleteResponse struct{}
func Delete(req DeleteRequest) (DeleteResponse, error) { ... }
```

#### `pkg/cache/set.go`, `pkg/cache/get.go`, `pkg/cache/delete.go`
session과 동일 패턴.

#### `pkg/file/upload.go`
```go
// @func upload
type UploadRequest struct { Key string, Body string }
type UploadResponse struct { Key string }
func Upload(req UploadRequest) (UploadResponse, error) { ... }
```

#### `pkg/file/download.go`, `pkg/file/delete.go`
```go
// @func download / delete
```

Func body는 내부적으로 `SessionModel`/`CacheModel`/`FileModel` 인터페이스를 호출한다. backend 선택은 `os.Getenv()`로 처리 (purity 규칙에서 `os` 허용).

### 7. crosscheck: 금지 import 목록 재정의

#### `internal/crosscheck/check_forbidden_imports.go`
모든 `@call` Func에 동일하게 적용한다. 면제 로직 없이, 금지 대상을 DB·네트워크 API로 한정한다.

**변경 전** (9개 금지):
```go
var forbiddenImportPrefixes = []string{
    "database/sql",
    "github.com/lib/pq",
    "github.com/jackc/pgx",
    "net/http",
    "net/rpc",
    "google.golang.org/grpc",
    "io",       // ← 삭제
    "io/ioutil", // ← 삭제
    "bufio",    // ← 삭제
}
```

**변경 후** (6개 금지):
```go
var forbiddenImportPrefixes = []string{
    "database/sql",
    "github.com/lib/pq",
    "github.com/jackc/pgx",
    "net/http",
    "net/rpc",
    "google.golang.org/grpc",
}
```

**원칙**: 모든 `@call` Func에서 DB 직접 접근과 API 호출은 절대 금지. file I/O(`io`, `bufio`, `os`)와 session/cache 읽기·쓰기는 모든 Func에서 허용.

### 8. 매뉴얼 수정 (`artifacts/manual-for-ai.md`)

#### 삭제
- L259~274 "Package-Prefix @model (Non-DDL Models)" 섹션 전체 삭제
- L397~426 "Built-in Models (pkg/)" 섹션에서 session, cache, file 제거

#### 추가/수정
- "Built-in Functions (pkg/)" 섹션에 session, cache, file 추가:

```
#### session
| Function | Request Fields | Response Fields | @error | Source |
|---|---|---|---|---|
| `set` | `Key`, `Value`, `TTL` | (none) | — | pkg/session |
| `get` | `Key` | `Value` | 404 | pkg/session |
| `delete` | `Key` | (none) | — | pkg/session |

Backend configured via `fullend.yaml` `session.backend` (postgres | memory).

#### cache
| Function | Request Fields | Response Fields | @error | Source |
|---|---|---|---|---|
| `set` | `Key`, `Value`, `TTL` | (none) | — | pkg/cache |
| `get` | `Key` | `Value` | — | pkg/cache |
| `delete` | `Key` | (none) | — | pkg/cache |

Backend configured via `fullend.yaml` `cache.backend` (postgres | memory).

#### file
| Function | Request Fields | Response Fields | @error | Source |
|---|---|---|---|---|
| `upload` | `Key`, `Body` | `Key` | — | pkg/file |
| `download` | `Key` | `Body` | 404 | pkg/file |
| `delete` | `Key` | (none) | — | pkg/file |

Backend configured via `fullend.yaml` `file.backend` (s3 | local).
```

- Purity Rule 수정:
```
All `@call func` rules:
- ALLOWED: file I/O (`io`, `bufio`, `os`), session/cache read/write
- FORBIDDEN: DB access (`database/sql`, `github.com/lib/pq`, `github.com/jackc/pgx`)
- FORBIDDEN: API calls (`net/http`, `net/rpc`, `google.golang.org/grpc`)
No per-package exceptions — same rule for all @call funcs.
```

- Cross-Validation Rules 추가:
```
| package-prefix @model used | ERROR |
```

- Fallback Chain 보충:
```
1. `specs/<project>/func/<pkg>/` — Project custom (overrides default)
2. `pkg/<pkg>/` — fullend default (session, cache, file, auth, etc.)
3. Neither → ERROR with skeleton suggestion
```

### 9. BUG029 보고서 수정

`files/bugs/BUG029.md` — 코드젠 버그에서 SSaC 설계 결함으로 재분류. 해결 방법: package-prefix @model 폐기.

### 10. filefunc validate 해소

기존 8건 위반 수정.

| 파일 | 위반 | 수정 방법 |
|---|---|---|
| `internal/crosscheck/check_input_key_case.go` | Q1 중첩 5 | 내부 루프를 헬퍼 함수로 추출 |
| `internal/crosscheck/check_jwt_builtin_inputs.go` | F1 함수 2개 + Q1 중첩 4 | 함수 분리 (파일 2개) + 내부 루프 추출 |
| `internal/crosscheck/collect_ddl_role_values.go` | Q1 중첩 3 | 내부 루프를 헬퍼 함수로 추출 |
| `internal/funcspec/collect_package_types.go` | Q1 중첩 4 | 내부 루프를 헬퍼 함수로 추출 |
| `internal/funcspec/fill_missing_fields.go` | Q1 중첩 3 | 내부 루프를 헬퍼 함수로 추출 |
| `internal/orchestrator/validate_funcspec.go` | Q1 중첩 3 | 내부 루프를 헬퍼 함수로 추출 |
| `internal/ssac/validator/err_ctx.go` | F3 메서드 2개 | 메서드 분리 (파일 2개) |

Phase051 신규/수정 파일도 filefunc 규칙 준수 확인.

최종: `filefunc validate` → 0 violation

---

## 변경 파일 요약

| 파일 | 변경 | 종류 |
|---|---|---|
| `internal/ssac/validator/validate_model.go` | package-prefix ERROR 차단 | 수정 |
| `internal/ssac/validator/validate_package_model.go` | 삭제 | 삭제 |
| `internal/ssac/validator/load_package_interfaces.go` | 삭제 | 삭제 |
| `internal/ssac/validator/parse_package_interfaces.go` | 삭제 | 삭제 |
| `internal/ssac/validator/load_package_go_interfaces.go` | 삭제 | 삭제 |
| `internal/ssac/validator/validator_test.go` | package-prefix 테스트 → ERROR 기대로 변경 | 수정 |
| `internal/ssac/generator/collect_model_usages_from_func.go` | `seq.Package` 조건 제거 | 수정 |
| `internal/ssac/generator/default_message.go` | `seq.Package` 참조 제거 | 수정 |
| `internal/ssac/generator/go_interface_test.go` | `TestGeneratePackageModelSkipInterface` 삭제 | 수정 |
| `internal/ssac/generator/go_handler_test.go` | `TestGeneratePackageModelCall` 삭제 | 수정 |
| `internal/crosscheck/check_forbidden_imports.go` | 금지 목록에서 io/bufio 제거 (DB·API만 금지) | 수정 |
| `pkg/session/set.go` | @func set + Request/Response | 신규 |
| `pkg/session/get.go` | @func get + Request/Response | 신규 |
| `pkg/session/delete.go` | @func delete + Request/Response | 신규 |
| `pkg/cache/set.go` | @func set + Request/Response | 신규 |
| `pkg/cache/get.go` | @func get + Request/Response | 신규 |
| `pkg/cache/delete.go` | @func delete + Request/Response | 신규 |
| `pkg/file/upload.go` | @func upload + Request/Response | 신규 |
| `pkg/file/download.go` | @func download + Request/Response | 신규 |
| `pkg/file/delete.go` | @func delete + Request/Response | 신규 |
| `artifacts/manual-for-ai.md` | package-prefix 삭제, session/cache/file Func 추가 | 수정 |
| `files/bugs/BUG029.md` | SSaC 설계 결함으로 재분류 | 신규 |

+ filefunc 해소 7파일

## 검증 방법

### 단위 테스트
- `go test ./...` 전체 통과
- `go build ./...` 전체 통과
- package-prefix SSaC 파싱 → validate ERROR 확인

### 더미 프로젝트 검증
zenflow-try06 add05(스케줄)를 `@call session.Set/Get/Delete`로 재작성:
1. `fullend validate` → ERROR 0
2. `fullend gen` → 코드 생성
3. `go build` → 컴파일 성공
4. `hurl --test` → scenario-schedule.hurl 통과

### filefunc validate
`filefunc validate` → 0 violation

## 의존성
Phase050 이후 독립 작업.
