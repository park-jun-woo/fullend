# Phase052 — BUG030 request body 판단 수정 + BUG031 session/cache/file Init 생성

## 목표
1. OpenAPI에 `requestBody`가 있는 엔드포인트는 `ShouldBindJSON`으로 JSON body를 읽도록 수정
2. `fullend.yaml`에 `session`/`cache`/`file` backend가 설정되면 `main.go`에 Init 호출을 생성

## 배경
zenflow-try06 add05(스케줄)에서 `@call`만 있는 POST 엔드포인트의 JSON body가 `c.Query()`로 생성되어 빈 문자열이 되는 버그(BUG030)와, `session.Init()`가 main.go에 미생성되어 nil panic이 발생하는 버그(BUG031) 발견.

---

## 변경 내용

### 1. BUG030 — OperationSymbol에 HasRequestBody 추가 + shouldUseJSONBody 수정

#### 원인 분석

`shouldUseJSONBody`가 `@post`/`@put` DDL 시퀀스 존재 여부로 JSON body 사용을 판단. SetSchedule처럼 `@call`만 있는 POST 엔드포인트는 `hasBodySeq = false`, rawParams 1개 → `c.Query()` 생성.

실제로는 **OpenAPI에 requestBody가 있느냐**가 기준이어야 함. GET 엔드포인트의 `request.*`는 query param이므로 `c.Query()` 유지, POST/PUT/DELETE의 requestBody 필드는 `ShouldBindJSON` 사용.

#### `internal/ssac/validator/operation_symbol.go`

`HasRequestBody bool` 필드 추가.

```go
type OperationSymbol struct {
    RequestFields  map[string]bool
    PathParams     []PathParam
    HasRequestBody bool  // OpenAPI requestBody 존재 여부
    XPagination    *XPagination
    XSort          *XSort
    XFilter        *XFilter
    XInclude       *XInclude
}
```

#### `internal/ssac/validator/symbol_table_build_operation_symbol.go`

`op.RequestBody != nil`이면 `opSym.HasRequestBody = true` 설정.

```go
if op.RequestBody == nil {
    return opSym
}
opSym.HasRequestBody = true
```

#### `internal/ssac/generator/should_use_json_body.go`

시그니처에 operationID 추가. SymbolTable에서 `HasRequestBody`를 조회.

**변경 전**:
```go
func shouldUseJSONBody(seqs []parser.Sequence, st *validator.SymbolTable, rawParams []rawParam) bool {
    hasBodySeq := false
    for _, seq := range seqs {
        if seq.Type == parser.SeqPost || seq.Type == parser.SeqPut {
            hasBodySeq = true
            break
        }
    }
    return (st != nil && len(rawParams) >= 2) || (hasBodySeq && len(rawParams) >= 1)
}
```

**변경 후**:
```go
func shouldUseJSONBody(seqs []parser.Sequence, st *validator.SymbolTable, operationID string, rawParams []rawParam) bool {
    if len(rawParams) == 0 {
        return false
    }
    // 1차: OpenAPI에 requestBody가 있으면 JSON body
    if st != nil {
        if op, ok := st.Operations[operationID]; ok {
            return op.HasRequestBody
        }
    }
    // 2차 fallback: Operations 미등록 시 @post/@put 시퀀스로 판단 (테스트 호환)
    for _, seq := range seqs {
        if seq.Type == parser.SeqPost || seq.Type == parser.SeqPut {
            return true
        }
    }
    return false
}
```

**근거**:
- 1차: OpenAPI `HasRequestBody`로 판단 — 가장 정확 (GET은 false, POST+body는 true)
- 2차 fallback: Operations가 빈 맵인 단위 테스트에서 `@post`/`@put` 시퀀스로 판단 — 기존 동작 호환
- rawParams가 0이면 바인딩 자체 불필요

#### `internal/ssac/generator/collect_request_params.go`

`shouldUseJSONBody` 호출부에 `operationID` 추가.

```go
if shouldUseJSONBody(seqs, st, operationID, rawParams) {
```

#### 테스트 영향

- `TestGenerateGet` (st=nil): rawParams 1개, st=nil → fallback → `@post` 없음 → `false` → `c.Query()`. **변화 없음**.
- `TestGeneratePost` (st=nil): rawParams 2개, st=nil → fallback → `@post` 있음 → `true` → `ShouldBindJSON`. **변화 없음**.
- `TestGenerateWithJSONBody` (st 있지만 Operations 빈맵): rawParams 2개, Operations miss → fallback → `@post` 있음 → `true` → `ShouldBindJSON`. **변화 없음**.
- `TestGenerateWithPathParam` (st 있고 Operations 등록): rawParams 0개(path param 제외) → `len(rawParams) == 0 → false`. **변화 없음**.
- 실제 fullend gen (Operations 등록됨): `HasRequestBody`로 정확히 판단. SetSchedule(POST+body) → `true`. **BUG030 해결**.

### 2. BUG031 — session/cache/file Init 생성

#### 설계 원칙

- session/cache는 `BuiltinBackend.Backend` 문자열(postgres|memory)만으로 Init 코드 생성 가능
- file은 backend별로 추가 설정(`LocalConfig.Root`, `S3Config.Bucket/Region`)이 필요하므로 `*projectconfig.FileBackend` 전체를 전달

#### `internal/gen/gogin/generate.go`

`parsed.Config`에서 Session/Cache backend 문자열과 File 설정 전체를 추출하여 `generateMainWithDomains`와 `generateMain`에 전달.

```go
var sessionBackend, cacheBackend string
var fileConfig *projectconfig.FileBackend
if parsed.Config != nil {
    if parsed.Config.Session != nil {
        sessionBackend = parsed.Config.Session.Backend
    }
    if parsed.Config.Cache != nil {
        cacheBackend = parsed.Config.Cache.Backend
    }
    fileConfig = parsed.Config.File  // nil이면 file Init 미생성
}
```

호출부 변경:

```go
// domain mode
generateMainWithDomains(cfg.ArtifactsDir, parsed.ServiceFuncs, cfg.ModulePath,
    queueBackend, parsed.Policies, sessionBackend, cacheBackend, fileConfig)

// flat mode
generateMain(cfg.ArtifactsDir, models, cfg.ModulePath,
    queueBackend, parsed.ServiceFuncs, parsed.Policies, sessionBackend, cacheBackend, fileConfig)
```

#### `internal/gen/gogin/generate_main_with_domains.go`

import에 `"strings"`, `"github.com/park-jun-woo/fullend/internal/projectconfig"` 추가.

```go
import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/park-jun-woo/fullend/internal/policy"
    "github.com/park-jun-woo/fullend/internal/projectconfig"
    ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)
```

시그니처에 `sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend` 추가.

```go
func generateMainWithDomains(artifactsDir string, serviceFuncs []ssacparser.ServiceFunc,
    modulePath string, queueBackend string, policies []*policy.Policy,
    sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend) error {
```

`buildQueueBlocks` 호출 아래에서 `buildBuiltinInitBlocks` 호출 추가:

```go
builtinImport, builtinInitBlock := buildBuiltinInitBlocks(sessionBackend, cacheBackend, fileConfig)
```

`mainWithDomainsTemplate` 호출부에 `builtinImport`, `builtinInitBlock` 인자 추가:

```go
src := mainWithDomainsTemplate(osImport, importBlock, queueImport, builtinImport,
    jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock)
```

#### `internal/gen/gogin/generate_main.go`

import에 `"github.com/park-jun-woo/fullend/internal/projectconfig"` 추가. (`"strings"`는 이미 존재.)

```go
import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/park-jun-woo/fullend/internal/policy"
    "github.com/park-jun-woo/fullend/internal/projectconfig"
    ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)
```

시그니처에 동일하게 `sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend` 추가.

```go
func generateMain(artifactsDir string, models []string, modulePath string,
    queueBackend string, serviceFuncs []ssacparser.ServiceFunc, policies []*policy.Policy,
    sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend) error {
```

`buildBuiltinInitBlocks` 호출 및 `mainTemplate` 호출부 변경 (domain mode와 동일).

#### `internal/gen/gogin/build_builtin_init_block.go` (신규)

session/cache/file Init 코드 블록 + import 빌더. `buildQueueBlocks` 패턴과 동일하게 (import, initBlock) 튜플 반환.

```go
package gogin

import (
    "fmt"
    "strings"

    "github.com/park-jun-woo/fullend/internal/projectconfig"
)

func buildBuiltinInitBlocks(sessionBackend, cacheBackend string, fileConfig *projectconfig.FileBackend) (builtinImport, builtinInitBlock string) {
    var imports []string
    var inits []string

    // --- session ---
    if sessionBackend == "postgres" {
        imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/session"`)
        inits = append(inits, `
	sm, err := session.NewPostgresSession(context.Background(), conn)
	if err != nil {
		log.Fatalf("session init failed: %v", err)
	}
	session.Init(sm)`)
    } else if sessionBackend == "memory" {
        imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/session"`)
        inits = append(inits, `
	session.Init(session.NewMemorySession())`)
    }

    // --- cache ---
    if cacheBackend == "postgres" {
        imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/cache"`)
        inits = append(inits, `
	cm, err := cache.NewPostgresCache(context.Background(), conn)
	if err != nil {
		log.Fatalf("cache init failed: %v", err)
	}
	cache.Init(cm)`)
    } else if cacheBackend == "memory" {
        imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/cache"`)
        inits = append(inits, `
	cache.Init(cache.NewMemoryCache())`)
    }

    // --- file ---
    if fileConfig != nil {
        imports = append(imports, `"github.com/park-jun-woo/fullend/pkg/file"`)
        switch fileConfig.Backend {
        case "local":
            root := "./uploads"
            if fileConfig.Local != nil && fileConfig.Local.Root != "" {
                root = fileConfig.Local.Root
            }
            inits = append(inits, fmt.Sprintf(`
	file.Init(file.NewLocalFile(%q))`, root))
        case "s3":
            bucket := ""
            region := "ap-northeast-2"
            if fileConfig.S3 != nil {
                bucket = fileConfig.S3.Bucket
                region = fileConfig.S3.Region
            }
            imports = append(imports, `"github.com/aws/aws-sdk-go-v2/config"`)
            imports = append(imports, `"github.com/aws/aws-sdk-go-v2/service/s3"`)
            inits = append(inits, fmt.Sprintf(`
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(%q))
	if err != nil {
		log.Fatalf("aws config failed: %%v", err)
	}
	file.Init(file.NewS3File(s3.NewFromConfig(awsCfg), %q))`, region, bucket))
        }
    }

    // --- context import (postgres session/cache에서 사용) ---
    needsContext := sessionBackend == "postgres" || cacheBackend == "postgres" ||
        (fileConfig != nil && fileConfig.Backend == "s3")
    if needsContext {
        hasContext := false
        for _, imp := range imports {
            if imp == `"context"` {
                hasContext = true
                break
            }
        }
        if !hasContext {
            imports = append([]string{`"context"`}, imports...)
        }
    }

    if len(imports) > 0 {
        builtinImport = "\n\t" + strings.Join(imports, "\n\t")
    }
    if len(inits) > 0 {
        builtinInitBlock = strings.Join(inits, "")
    }
    return builtinImport, builtinInitBlock
}
```

**설계 결정 — file backend 인자:**

| backend | 인자 출처 | 기본값 |
|---|---|---|
| `file/local` | `fullend.yaml` → `file.local.root` | `"./uploads"` |
| `file/s3` | `fullend.yaml` → `file.s3.bucket`, `file.s3.region` | region: `"ap-northeast-2"` |

- `file/local`: `LocalConfig.Root`가 비어있으면 `"./uploads"` 기본값 사용. root 경로를 하드코딩하되 `fullend.yaml`에서 오버라이드 가능.
- `file/s3`: AWS SDK `config.LoadDefaultConfig`로 클라이언트 생성. `S3Config.Bucket`/`Region`은 `fullend.yaml`에서 읽음. 크레덴셜은 AWS 표준 환경변수(`AWS_ACCESS_KEY_ID` 등) 위임.

**context import 중복 방지:** queue가 이미 `"context"`를 import하는 경우가 있으므로, template에서 `%s` 삽입 시 중복 import가 되지 않도록 주의 필요. 단, 현재 구조상 각 블록이 자체 import 문자열을 반환하고 template의 `import (...)` 블록에 문자열로 합쳐지므로, **queue가 `context`를 import할 때 builtin 쪽에서는 `context`를 빼야 함**. → `buildBuiltinInitBlocks`에 `queueImport string` 인자를 추가하거나, template 레벨에서 중복 제거하는 방식이 필요.

**채택 방안:** `buildBuiltinInitBlocks`가 `context`를 import 목록에 포함하되, 호출 측(`generateMainWithDomains`/`generateMain`)에서 queueImport에 `"context"`가 포함되어 있으면 builtinImport에서 `"context"`를 제거한다.

```go
// generateMainWithDomains / generateMain에서:
builtinImport, builtinInitBlock := buildBuiltinInitBlocks(sessionBackend, cacheBackend, fileConfig)
if strings.Contains(queueImport, `"context"`) {
    builtinImport = strings.Replace(builtinImport, "\n\t\"context\"", "", 1)
}
```

#### `internal/gen/gogin/main_with_domains_template.go`

`builtinImport`, `builtinInitBlock` 인자 2개 추가 (총 10개 인자).

**변경 전**:
```go
func mainWithDomainsTemplate(osImport, importBlock, queueImport,
    jwtFlagLine, authzBlock, queueInitBlock, initBlock, queueSubscribeBlock string) string {
```

**변경 후**:
```go
func mainWithDomainsTemplate(osImport, importBlock, queueImport, builtinImport,
    jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock string) string {
```

template 본문 import 섹션:
```go
	_ "github.com/lib/pq"
%s%s%s
)
```
→ `importBlock`, `queueImport`, `builtinImport` 순으로 3개 `%s` 삽입.

template 본문 init 섹션:
```go
%s%s%s
	server := &service.Server{
%s
	}
%s
```
→ `authzBlock`, `queueInitBlock`, `builtinInitBlock`, `initBlock`, `queueSubscribeBlock` 순으로 5개 `%s` 삽입.

#### `internal/gen/gogin/main_template.go`

동일하게 `builtinImport`, `builtinInitBlock` 인자 2개 추가 (총 9개 인자).

**변경 전**:
```go
func mainTemplate(modulePath, authzImport, queueImport,
    authzInitBlock, queueInitBlock, initBlock, queueSubscribeBlock string) string {
```

**변경 후**:
```go
func mainTemplate(modulePath, authzImport, queueImport, builtinImport,
    authzInitBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock string) string {
```

template 본문 import 섹션:
```go
	"%s/internal/service"%s%s%s
)
```
→ `modulePath`, `authzImport`, `queueImport`, `builtinImport` 순으로 4개 `%s` 삽입.

template 본문 init 섹션:
```go
%s%s%s
	server := &service.Server{
%s
	}
%s
```
→ `authzInitBlock`, `queueInitBlock`, `builtinInitBlock`, `initBlock`, `queueSubscribeBlock` 순으로 5개 `%s` 삽입.

---

## 변경 파일 요약

| 파일 | 변경 | 종류 |
|---|---|---|
| `internal/ssac/validator/operation_symbol.go` | `HasRequestBody bool` 필드 추가 | 수정 |
| `internal/ssac/validator/symbol_table_build_operation_symbol.go` | `HasRequestBody = true` 설정 | 수정 |
| `internal/ssac/generator/should_use_json_body.go` | `OperationSymbol.HasRequestBody`로 판단 | 수정 |
| `internal/ssac/generator/collect_request_params.go` | `shouldUseJSONBody` 호출부 수정 | 수정 |
| `internal/gen/gogin/generate.go` | session/cache backend 문자열 + `*FileBackend` 추출·전달 | 수정 |
| `internal/gen/gogin/generate_main_with_domains.go` | 시그니처 + `buildBuiltinInitBlocks` 호출 + context 중복 제거 | 수정 |
| `internal/gen/gogin/generate_main.go` | 시그니처 + `buildBuiltinInitBlocks` 호출 + context 중복 제거 | 수정 |
| `internal/gen/gogin/build_builtin_init_block.go` | Init 코드 블록 + import 빌더 | 신규 |
| `internal/gen/gogin/main_with_domains_template.go` | `builtinImport` + `builtinInitBlock` 인자 추가 (8→10개) | 수정 |
| `internal/gen/gogin/main_template.go` | `builtinImport` + `builtinInitBlock` 인자 추가 (7→9개) | 수정 |

## 검증 방법

### 단위 테스트
- `go test ./...` 전체 통과
- `go build ./...` 전체 통과

### 더미 프로젝트 검증
zenflow-try06 (`session.backend: postgres`, `file.backend: local`):
1. `fullend validate` → ERROR 0
2. `fullend gen` → 코드 생성
3. SetSchedule 핸들러에 `ShouldBindJSON` 생성 확인
4. main.go에 `session.Init(session.NewPostgresSession(...))` 생성 확인
5. main.go에 `file.Init(file.NewLocalFile("./uploads"))` 생성 확인
6. main.go import에 `"github.com/park-jun-woo/fullend/pkg/session"` 존재 확인
7. main.go import에 `"github.com/park-jun-woo/fullend/pkg/file"` 존재 확인
8. context import 중복 없음 확인
9. `go build` → 컴파일 성공
10. `hurl --test scenario-schedule.hurl` → 통과

## 의존성
Phase051 이후. 독립 작업.
