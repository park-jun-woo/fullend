# internal/gen/gogin

OpenAPI + SSaC + DDL + fullend.yaml + Policy + StateDiagram을 조합해 **Go + Gin 백엔드**를 생성하는 코드젠 엔진.
132개 파일, 약 4,954 LOC. `fullend`의 가장 크고 오래된 생성기.

## 용어 주석

본 README에서 **feature**는 `specs/service/<폴더>/*.ssac`의 서브폴더 이름이다. 내부 코드에서는 `Domain`이라는 이름을 사용한다 (`ServiceFunc.Domain`, `hasDomains()`, `generateServerStructWithDomains()` 등). 두 용어는 동의어로 읽는다.

## 진입점

```go
// internal/gen/gogin/generate.go
func (g *GoGin) Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig) error
```

최상위 분기 (`generate.go:51-92`):

```
hasDomains(ServiceFuncs)   ─┐
                            ├─ true  → Feature mode 호출 계열 (…WithDomains)
                            └─ false → Flat mode 호출 계열
```

**Flat mode**는 단일 `Server` 구조체에 모든 모델·함수 필드를 집약하고 `http.ServeMux`를 사용한다.
**Feature mode**는 feature별 `Handler`를 두고 중앙 `Server`가 이들을 필드로 합쳐 `gin.Engine`에 연결한다.

파서(`ssac/parser/parse_dir_entry.go:19`)가 `service/` 루트 직접 파일을 금지하므로 실전에서는 **항상 Feature mode**로 진입한다. Flat 경로는 역사적 잔존.

## 산출물 맵

### Feature mode (실전 경로)

| 경로 | 책임 함수 | 내용 |
|------|----------|------|
| `backend/cmd/main.go` | `generateMainWithDomains()` `generate_main_with_domains.go:18` | DB 연결, feature Handler 초기화, session/cache/file/queue/authz init, `gin.Run()` |
| `backend/internal/service/server.go` | `generateCentralServer()` `generate_central_server.go:18` | 중앙 `Server` (feature Handler 필드 집합) + `SetupRouter()` |
| `backend/internal/service/{feature}/handler.go` | `generateDomainHandler()` `generate_domain_handler.go:16` | feature별 `Handler` 구조 (DB, Model, @call 함수, JWTSecret 필드) |
| `backend/internal/service/{feature}/*.go` | `transformServiceFilesWithDomains()` `transform_service_files_with_domains.go:17` | SSaC가 생성한 서비스 함수를 `(h *Handler)` 메서드로 변환 |
| `backend/internal/model/types.go` | `generateTypesFile()` `generate_types_file.go` | DDL에서 파생된 struct 타입 |
| `backend/internal/model/{model}.go` | `generateModelFile()` `generate_model_file.go:15` | `{Model}ModelImpl` + 메서드(Create/FindByID/List/Update/Delete/WithTx) |
| `backend/internal/model/queryopts.go` | `generateQueryOpts()` `generate_query_opts.go:12` | `ParseQueryOpts`, `BuildSelectQuery`, `BuildCountQuery` |
| `backend/internal/model/include_helpers.go` | `generateIncludeHelpersFile()` `generate_model_file.go:72` | FK 관계 로딩 헬퍼 |
| `backend/internal/model/auth.go` | `generateAuthStubWithDomains()` `generate_auth_stub_with_domains.go:17` | `CurrentUser` (claims 파생), `Authorizer` interface |
| `backend/internal/auth/{issue,verify,refresh}_token.go` | `generateAuthPackage()` `generate_auth_package.go:14` | JWT 발급/검증/갱신 |
| `backend/internal/middleware/bearerauth.go` | `generateMiddleware()` `generate_middleware.go:17` | BearerAuth gin middleware |
| `backend/internal/authz/*.rego` | `GenerateAuthzPackage()` `authz.go:14` | Rego 정책 파일 복사 + `authz.OwnershipMapping` 리터럴 |
| `backend/internal/states/{id}state/{id}state.go` | `generateSingleStateMachine()` `generate_state_machines.go:14` | 상태 전이 머신 |

### Flat mode 차이점 (참고용)

| 변경점 | 함수 |
|-------|------|
| `Server` 하나에 모든 필드 | `generateServerStruct()` `generate_server_struct.go:16` |
| `(s *Server)` receiver 변환 | `transformServiceFiles()` `transform_service_files.go:18` |
| `http.ServeMux` 라우팅 | `generate_server_struct.go:53-70` |
| `backend/internal/service/auth.go` 고정 `CurrentUser` | `generateAuthStub()` `auth.go:12` |

## 입력 의존성

`ParsedSSOTs`에서 소비하는 필드 (generate.go 상단부):

| 필드 | 용도 |
|------|------|
| `ServiceFuncs[].Name` | operationId, 메서드명 |
| `ServiceFuncs[].Domain` | feature mode 분기 및 디렉토리 배치 |
| `ServiceFuncs[].Sequences` | seqType, 사용 모델, @call 함수, subscribe/publish |
| `ServiceFuncs[].Param` | subscriber 함수 파라미터 타입 |
| `OpenAPIDoc.Paths` | 라우트 경로/메서드 |
| `OpenAPIDoc.Operations` | OperationID, Security, x-include, x-pagination, x-sort, PathParameters |
| `OpenAPIDoc.Components.SecuritySchemes` | bearerAuth 탐지 |
| `Policies[].Ownerships` | `authz.Init()` 리터럴 |
| `Config.Backend.Auth.Claims` | JWT claim 필드 정의 |
| `Config.Backend.Auth.SecretEnv` | JWT secret 환경변수 이름 |
| `Config.Queue.Backend` | queue init |
| `Config.Session.Backend` | session init |
| `Config.Cache.Backend` | cache init |
| `Config.File` | file init (S3/local) |

또한 `cfg.SpecsDir`에서 직접 파싱하는 것:
- DDL 파일 → 테이블/컬럼/FK 메타 (`parse_ddl_files.go`)
- sqlc query 파일 → SQL 문자열 및 cardinality (`parse_query_files.go`)
- 이미 존재하는 `models_gen.go` → 인터페이스 시그니처 (`parse_models_gen.go:`)

## 파이프라인 (12단계)

1. **메타 수집** — `collectModels`, `collectFuncs`, `hasDomains`, `uniqueDomains`.
2. **구조 검증** — `hasBearerScheme(doc)` → claims 필수 여부.
3. **서비스 파일 변환** — `transformServiceFiles[WithDomains]()`:
   - `func Xxx(...)` → `func (r *Receiver) Xxx(...)` (r은 Server 또는 Handler)
   - 모델 참조: `courseModel.X` → `r.CourseModel.X`
   - @call 참조: `hashPassword(` → `r.HashPassword(`
   - `__RESPONSE_STATUS__` → OpenAPI success code
4. **Server/Handler 구조 생성** — `generateServerStruct[WithDomains]()` + `generateDomainHandler()`.
5. **main.go 생성** — `generateMain[WithDomains]()` + 초기화 블록 조합.
6. **모델 인터페이스 파싱** — `parseModelsGen()`: `models_gen.go`에서 `{Model}Model` 시그니처 추출.
7. **DDL/Query 파싱** — `parseDDLFiles()`, `parseQueryFiles()`, `collectSeqTypes()`.
8. **모델 구현 생성** — `generateModelFile()` 각 모델마다: constructor + `scanX` + CRUD 메서드.
9. **Pagination 구현** — `writeOffsetPaginationMethod()` / `writeCursorPaginationMethod()` + `generateQueryOpts()`.
10. **Include(FK) 처리** — `collectModelIncludes()`, `resolveIncludes()`, `generateIncludeHelpersFile()`.
11. **Auth/Authz 생성** — `generateAuthPackage()`, `generateMiddleware()`, `GenerateAuthzPackage()`.
12. **디렉티브 주입** — `attachServiceDirectives()`, `attachTSXDirectives()` (프론트엔드 쪽 파일에도 마킹).

## Model 생성 스펙

### 인터페이스 읽기 → 구현 쓰기

외부 도구(SSaC generator)가 먼저 `internal/model/models_gen.go`에 인터페이스를 써둠. gogin은 그걸 **파싱**해서 구현체를 채우는 구조 (`parseModelsGen()` `parse_models_gen.go:14`).

### 메서드 라우팅 규칙

`generateMethodFromIface()` `generate_method_from_iface.go:12-74` 가 메서드 이름 + 반환 타입 + seqType 힌트 조합으로 구현을 선택한다.

| 조건 | 구현 함수 | 반환 타입 예시 |
|------|----------|---------------|
| 이름이 `WithTx` | 인라인 특수 케이스 | `{Model}Model` |
| List + QueryOpts + `pagination.Cursor[T]` 반환 | `writeCursorPaginationMethod` | `(*pagination.Cursor[T], error)` |
| List + QueryOpts + `pagination.Page[T]` 반환 | `writeOffsetPaginationMethod(isPageReturn=true)` | `(*pagination.Page[T], error)` |
| List + QueryOpts | `writeOffsetPaginationMethod(isPageReturn=false)` | `([]T, int64, error)` |
| `[]` 반환 (slice) | `writeSliceReturnMethod` | `([]T, error)` |
| `Find…` 접두 또는 seqType=`get` | `writeFindMethod` | `(*{Model}, error)` |
| seqType=`post` | 인라인 `QueryRowContext` + `scanX` | `(*{Model}, error)` |
| seqType=`put`/`delete` | 인라인 `ExecContext` | `error` |
| 기타 + `query.Cardinality=="one"` | 인라인 `QueryRowContext` | `(*{Model}, error)` |
| 기타 + cardinality=`many` | 인라인 `QueryContext` + 루프 | `([]{Model}, error)` |

이 switch는 **축 5개**(이름 접두, QueryOpts 유무, 반환 제네릭, 슬라이스 여부, seqType) × **7 case**로 물려 있어, 케이스 추가 시 실수 여지 있음.

### Pagination 방식

**Offset**
```go
opts QueryOpts   // Limit, Offset, SortCol, SortDir, Filters
BuildCountQuery(table, baseWhere, baseArgCount, opts)
BuildSelectQuery(...)
// 반환: ([]T, int64, error) or (*pagination.Page[T], error)
```

**Cursor**
```go
opts QueryOpts   // Cursor, Limit, SortCol, SortDir
BuildSelectQuery(...)   // WHERE {sortCol} > cursor LIMIT {limit+1}
hasNext := len(items) > requestedLimit
nextCursor := fmt.Sprintf("%v", items[lastIdx].{CursorField})
```

### Include (FK 로딩)

`collectModelIncludes()` `collect_model_includes.go:17` 가 OpenAPI operation의 `x-include.allowed` 목록을 수집 → operationId → 모델 매핑.

`resolveIncludes()` 에서 DDL FK 정보와 합쳐 다음 구조 생성:

```go
includeMapping{
    IncludeName:  "instructor",       // "instructor_id"에서 _id 제거
    FieldName:    "Instructor",
    FieldType:    "*User",
    FKColumn:     "instructor_id",
    TargetModel:  "User",
}
```

List 메서드 내부에서 각 include에 대해 `m.include{Name}(items)` 호출이 생성됨.

### @sensitive

DDL 컬럼 주석에 `@sensitive`가 있으면 (`ddl_column.go`) struct field에 `json:"-"` 태그 부여.

## Handler 구조 (feature mode)

```go
// backend/internal/service/{feature}/handler.go
package {feature}

type Handler struct {
    DB         *sql.DB            // 이 feature의 seq에 post/put/delete 존재 시
    {Model1}Model  model.{Model1}Model
    {Model2}Model  model.{Model2}Model
    {Func1}    func(args ...interface{}) (interface{}, error)  // @call용 커스텀 함수
    JWTSecret  string             // 이 feature에 auth.IssueToken 존재 시
}
```

필드 선별 함수:
- `collectModelsForDomain()` — 이 feature가 실제 참조하는 모델만.
- `collectFuncsForDomain()` — 이 feature가 실제 참조하는 @call 함수만.
- `domainNeedsDB()` `domain_needs_db.go:9` — write seq 존재 여부.
- `domainNeedsJWTSecret()` `domain_needs_jwt_secret.go:9` — `auth.IssueToken` 존재 여부.

## 중앙 Server (feature mode)

```go
// backend/internal/service/server.go
type Server struct {
    Auth       *authsvc.Handler     // feature Handler 참조
    Gig        *gigsvc.Handler
    // ...
    UserModel  model.UserModel      // flat 잔존 (Domain="" 서비스함수용)
    JWTSecret  string               // hasBearerScheme(doc) 시
}

func SetupRouter(s *Server) *gin.Engine {
    r := gin.Default()
    auth := r.Group("/")
    auth.Use(middleware.BearerAuth(s.JWTSecret))

    auth.Handle("GET", "/gigs/:gigID", s.Gig.GetGig)
    r.Handle("POST", "/login", s.Auth.Login)
    return r
}
```

라우트 연결 (`write_central_routes.go:14`):
- OpenAPI Path × Method 순회
- `opDomains[operationID]` 조회로 feature 매핑
- `op.Security` 유무로 auth group vs 루트 분기
- `convertPathParamsGin()` — `{GigID}` → `:gigID` 변환

## main.go 초기화 블록

조건부로 켜지는 블록들:

| 블록 | 활성화 조건 | 생성 함수 |
|------|------------|----------|
| `authz.Init(conn, ownerships)` | `hasAuthSequence(ServiceFuncs)` | `build_ownerships_literal.go` |
| `queue.Init(...) + defer Close()` | `queueBackend != "" && (subscribers > 0 \|\| hasPublishSequence())` | `build_queue_blocks.go` |
| `queue.Subscribe(...)` × n | subscriber seq 수 | 동일 |
| `session.New{Pg,Mem}()` | `Config.Session.Backend` | `build_builtin_init_block.go` |
| `cache.New{Pg,Mem}()` | `Config.Cache.Backend` | 동일 |
| `file.NewS3() / NewLocal()` | `Config.File` | `build_file_init_block.go` |

Feature mode에서 추가:
- `flag.String("jwt-secret", os.Getenv("JWT_SECRET"), …)` — `anyDomainNeedsAuth()` 시 (`generate_main_with_domains.go`)
- 각 feature Handler `{Model, DB, JWTSecret}` 초기화 블록 — `build_domain_init_block.go`

## Auth / Authz / Middleware

### Auth 패키지

`generateAuthPackage()` `generate_auth_package.go:14` — `len(claims) > 0` 일 때만 생성.

| 파일 | 내용 |
|------|------|
| `internal/auth/issue_token.go` | `IssueTokenRequest` (claims 필드) + `IssueToken()` (access + exp 24h) |
| `internal/auth/verify_token.go` | `VerifyToken(token, secret)` → claims 추출 |
| `internal/auth/refresh_token.go` | `RefreshToken()` → 새 토큰 |
| `internal/auth/auth.go` | 위 함수 재export |

claims 필드 매핑:
- `claims[fieldName].Key` → JWT claim 키
- `claims[fieldName].GoType` → Go 타입
- 정렬: `sortedClaimFields()` (일관성 목적)

### Middleware

`generateMiddleware()` `generate_middleware.go:17` → `internal/middleware/bearerauth.go`:

```go
func BearerAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authorization 헤더 파싱
        // auth.VerifyToken()
        // c.Set("currentUser", &model.CurrentUser{...})
    }
}
```

### Authz

`GenerateAuthzPackage()` `authz.go:14` — `len(policies) > 0` 일 때.
- Rego 파일을 **생성이 아닌 복사** (`specs/authz/*.rego` → `backend/internal/authz/`).
- `authz.OwnershipMapping` 리터럴을 `build_ownerships_literal.go`가 Go 소스로 조립.

## Route 연결

### Flat mode (`generate_server_struct.go:53-70`)

```go
mux.HandleFunc("GET /gigs/{GigID}", withParams("GigID", s.GetGig))
```

- `convertPathParams()` — `{GigID}` 그대로 유지 (Go 1.22 mux 호환).
- `writeRouteHandler()` — 파라미터 있으면 추출 래퍼 감쌈.

### Feature mode (`write_central_routes.go:14`)

```go
auth.Handle("GET", "/gigs/:gigID", s.Gig.GetGig)
```

- `convertPathParamsGin()` — `{GigID}` → `:gigID` (lcFirst).
- auth group: `opHasSecurity(op) && hasBearer` 시.
- 타깃: flat은 `s.{op}`, feature는 `s.{Feature}.{op}`.

## Transform 규칙 상세

`transformSource()` `transform_source.go:15` 가 한 파일에 가하는 치환:

1. **함수 → 메서드** — `func X(...)` → `func (r *R) X(...)`.
2. **모델 참조** — `courseModel.X(...)` → `r.CourseModel.X(...)`.
3. **@call 참조** — `hashPassword(...)` → `r.HashPassword(...)` (단, Server/Handler에 필드 있는 것만).
4. **패키지 @call** — `auth.VerifyPassword(...)` 는 그대로 (import로 처리).
5. **status placeholder** — `__RESPONSE_STATUS__` → OpenAPI success code.
6. **imports 정리** — `fixImports()`.
7. **type assertions** — `addTypeAssertions()` (@call 반환 `interface{}`를 실제 타입으로 단언).

## 특이 패턴

### Directive 해시

생성된 파일 상단에 다음 디렉티브 삽입:

```go
//fullend:gen ownership="gen" ssot="fullend.yaml" contract={hash}
```

- `ownership="gen"` 기본값, `"preserve"`로 바꾸면 재생성 시 body 보존.
- `contract={hash}` — 입력 스펙 해시. 재생성 필요 여부 판단용 (예: `HashClaimDefs(claims)`).

### @call 함수 필드 vs 패키지 함수

| 시퀀스의 Model 문자열 | 처리 |
|----------------------|------|
| `hashPassword` (패키지 프리픽스 없음) | Server/Handler 필드로 추가 + receiver 참조로 치환 |
| `auth.VerifyPassword` (프리픽스 있음) | 필드 아님, import 경유 직접 호출 |

### JWT Secret 전달

- **Flat** — `os.Getenv(secretEnv)` 또는 하드코드 `"secret"`.
- **Feature** — `flag.String("jwt-secret", os.Getenv("JWT_SECRET"), …)` + 각 Handler에 주입.

### DDL 네이밍 변환

```go
// FOREIGN KEY instructor_id REFERENCES users(id)
ddlColumn.Name    = "instructor_id"   // snake_case
ddlColumn.GoName  = "InstructorID"    // PascalCase
ddlColumn.FKTable = "users"
```

### 명명 유틸

| 함수 | 예시 |
|------|------|
| `lcFirst("CourseModel")` | `"courseModel"` |
| `ucFirst("course")` | `"Course"` |
| `singularize("courses")` | `"Course"` (inflection + strcase) |

## 알려진 구조적 냄새 (리팩토링 참고)

1. **시그니처 비대** — `generateMain()` 8 파라미터, `generateMethodFromIface()` 9 파라미터. 결정 축이 시그니처에 산재.
2. **결정 분산** — 예: "queue init이 필요한가?"의 답이 `generate_main.go:68-97` + `collect_subscribers.go` + `has_publish_sequence.go` + `build_queue_blocks.go`로 분산.
3. **Flat vs Feature 중복** — 거의 모든 `generateXWithDomains()` 가 `generateX()` 와 7~80% 중복.
4. **7-case switch** — `generateMethodFromIface.go:36-73` 은 축 추가 시 실수 여지.
5. **Template + 동적 조립 혼합** — `main_template.go`, `main_with_domains_template.go`, `query_opts_template.go`는 큰 문자열 리터럴. 조건부 블록 삽입 지점이 문자열 치환이라 경계 취약.

위 냄새는 Toulmin 적용 또는 `ground.Build()` 기반 재구성의 가치를 판단하는 기준점.
