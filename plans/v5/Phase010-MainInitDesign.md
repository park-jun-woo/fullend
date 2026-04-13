# Phase010 Part B2 — MainInitDesign

> `pkg/generate/gogin/generate_main.go` 의 6축 초기화 블록 조합을 **`InitNeeds` struct + `DecideMainInit` 순수 함수**로 수렴하는 설계.
> 기준점 (2-depth) 재평가 반영 — Toulmin graph 불채택.

---

## 1. 현행 구조 분석

### 1.1 대상 함수

`generateMain(in MainGenInput) error` — `pkg/generate/gogin/generate_main.go:30~101`

`MainGenInput` (Phase009 수렴):

```go
type MainGenInput struct {
    ArtifactsDir   string
    ServiceFuncs   []ssacparser.ServiceFunc
    ModulePath     string
    QueueBackend   string                  // "", "redis", ...
    Policies       []*policy.Policy
    SessionBackend string                  // "postgres"/"memory"/""
    CacheBackend   string                  // "postgres"/"memory"/""
    FileConfig     *manifest.FileBackend
}
```

### 1.2 축별 판정 주체

| 축 | 판정 | 렌더링 | 코드 위치 |
|----|------|------|---------|
| AUTH | `anyDomainNeedsAuth(serviceFuncs, domains)` | authz.Init, os.Setenv JWT | L60-75 |
| AUTHZ | `anyNeedsAuth == true` (AUTH 와 동치) | policies 주입 | L67-75 |
| QUEUE | `queueBackend != "" ∧ (hasSubscribers ∨ hasPublishSeq)` | queue.Init + Subscribe + Start | `buildQueueBlocks` |
| SESSION | `sessionBackend ∈ {"postgres","memory"}` | session.NewX + Init | `buildBuiltinInitBlocks` |
| CACHE | `cacheBackend ∈ {"postgres","memory"}` | cache.NewX + Init | 동일 |
| FILE | `fileConfig != nil` (+ backend "local"/"s3") | file.NewLocal/S3 + Init | 동일 |

### 1.3 의존성

- **AUTH ↔ AUTHZ**: 1:1 동치 (공용 `anyNeedsAuth`). struct 에서는 분리 필드로 유지하되 판정값은 공유.
- **context import**: SESSION=postgres ∨ CACHE=postgres ∨ FILE=s3 중 1개 이상 활성 시 `"context"` import 필요. 현행 `generate_main.go:L81-83` 에서 수동 dedup.

### 1.4 depth 평가 (기준점 적용)

6축 판정을 순차 if 로 재표현:

```go
needs := InitNeeds{}
if anyDomainNeedsAuth(funcs, domains) {         // depth 1
    needs.Auth  = true
    needs.Authz = true
}
if queueBackend != "" && hasQueueWork(...) {    // depth 1
    needs.Queue = true
}
if sessionBackend != "" {                       // depth 1
    needs.Session = SessionInit{Enabled:true, Backend:sessionBackend}
}
if cacheBackend != "" {                         // depth 1
    needs.Cache = CacheInit{Enabled:true, Backend:cacheBackend}
}
if fileConfig != nil {                          // depth 1
    needs.File = FileInit{Enabled:true, Backend:fileConfig.Kind, Config:fileConfig}
}
needs.NeedsContextImport =
    (needs.Session.Enabled && needs.Session.Backend == "postgres") ||
    (needs.Cache.Enabled   && needs.Cache.Backend   == "postgres") ||
    (needs.File.Enabled    && needs.File.Backend    == "s3")
```

**최대 depth = 1**. **기준점 이내** → **Toulmin 제외**.

---

## 2. 설계

### 2.1 반환 타입

```go
// pkg/generate/gogin/decide_main_init.go
type InitNeeds struct {
    Auth    bool
    Authz   bool                      // AUTH 와 동치 (미래 분리 대비)
    Queue   bool

    Session SessionInit
    Cache   CacheInit
    File    FileInit

    NeedsContextImport bool           // 파생
}

type SessionInit struct { Enabled bool; Backend string }
type CacheInit   struct { Enabled bool; Backend string }
type FileInit    struct {
    Enabled bool
    Backend string
    Config  *manifest.FileBackend
}

type MainFacts struct {
    ServiceFuncs   []ssacparser.ServiceFunc
    Domains        []string
    QueueBackend   string
    HasSubscribers bool
    HasPublishSeq  bool
    SessionBackend string
    CacheBackend   string
    FileConfig     *manifest.FileBackend
}
```

### 2.2 판정 함수

```go
func DecideMainInit(facts MainFacts) InitNeeds {
    needs := InitNeeds{}

    if anyDomainNeedsAuth(facts.ServiceFuncs, facts.Domains) {
        needs.Auth  = true
        needs.Authz = true
    }
    if facts.QueueBackend != "" && (facts.HasSubscribers || facts.HasPublishSeq) {
        needs.Queue = true
    }
    if facts.SessionBackend != "" {
        needs.Session = SessionInit{Enabled: true, Backend: facts.SessionBackend}
    }
    if facts.CacheBackend != "" {
        needs.Cache = CacheInit{Enabled: true, Backend: facts.CacheBackend}
    }
    if facts.FileConfig != nil {
        needs.File = FileInit{
            Enabled: true,
            Backend: facts.FileConfig.Kind,
            Config:  facts.FileConfig,
        }
    }

    needs.NeedsContextImport =
        (needs.Session.Enabled && needs.Session.Backend == "postgres") ||
        (needs.Cache.Enabled   && needs.Cache.Backend   == "postgres") ||
        (needs.File.Enabled    && needs.File.Backend    == "s3")

    return needs
}
```

### 2.3 파일 배치

```
pkg/generate/gogin/
├── decide_main_init.go              신설 — InitNeeds/MainFacts 타입 + DecideMainInit
├── decide_main_init_test.go         신설 — 축별 on/off 조합 + 파생 ContextImport 테이블 테스트
└── generate_main.go                 수정 — 판정 제거, needs 소비
```

`anyDomainNeedsAuth` 등 기존 헬퍼는 그대로 사용.

### 2.4 호출자 소비 패턴

```go
func generateMain(in MainGenInput) error {
    facts := newMainFacts(in)
    needs := DecideMainInit(facts)

    var authzBlock string
    if needs.Authz {
        authzBlock = buildAuthzBlock(in.Policies, facts.Domains)
    }

    var queueImport, queueInit, queueSubs, queueStart string
    if needs.Queue {
        queueImport, queueInit, queueSubs, queueStart = buildQueueBlocks(in)
    }

    builtinImports, builtinInits := buildBuiltinInitBlocks(needs)   // InitNeeds 소비
    importBlock := mergeImports(builtinImports, queueImport, needs.NeedsContextImport)

    return mainWithDomainsTemplate(in.ArtifactsDir, importBlock,
        authzBlock, builtinInits, queueInit, queueSubs, queueStart, facts.Domains)
}
```

**변화**:
- `anyDomainNeedsAuth` 직접 호출 제거 → `needs.Auth` 참조
- `buildBuiltinInitBlocks` 시그니처가 `(needs InitNeeds)` 단일 인자로 수렴
- `mergeImports` 가 `needs.NeedsContextImport` bool 을 받아 수동 dedup 제거 (L81-83)

---

## 3. 검증

- `go test ./pkg/generate/gogin/...` — 6축 조합 테이블 (축별 on/off + 파생 ContextImport 8 case)
- `fullend gen dummys/gigbridge/specs /tmp/x` — 생성된 `cmd/main.go` 가 기존과 실질 동일 (import 순서·공백 차이 허용)
- `cd /tmp/x/backend && go build ./...` 성공

---

## 4. 보류

- `buildQueueBlocks` 내부의 `needsQueue := queueBackend != "" && ...` 중복 판정 제거. 호출은 `needs.Queue == true` 일 때만 발생하도록 계약 변경.
- `main_template.go` 의 `text/template` 화는 Phase010 범위 밖.
- 장래 AUTH/AUTHZ 분리(OIDC: auth 외부, authz 내부)가 필요해지면 `needs.Auth != needs.Authz` 조건 지원 가능.
- AUTH↔AUTHZ 분리 또는 축이 10+ 로 늘어나 **의존 그래프** 가 생기면 그때 Toulmin 승격 재평가.
