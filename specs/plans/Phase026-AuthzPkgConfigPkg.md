# ✅ Phase 026: pkg/authz + pkg/config + authzgen 전환 + config import 변환

## 목표

SSaC 수정지시서022 완료에 따른 fullend 측 연동 변경.

- `pkg/authz` — 기본 OPA 기반 인가 패키지 (CheckRequest/CheckResponse + Check 함수)
- `pkg/config` — 환경변수 기반 설정 조회 패키지
- `authzgen` 전환 — 고정 `Input{ID interface{}}` 자동생성 → `pkg/authz` 복사 + .rego 복사
- `gluegen` — SSaC가 생성한 `"config"` import → `pkg/config` 경로 변환
- `crosscheck` — `@auth` inputs → authz CheckRequest 필드 매칭 검증
- `fullend.yaml` — `authz.package` 설정 필드 추가

## 의존성

- SSaC 수정지시서 022 ✅ (A. @auth → @call 방식, B. @call 타입 검증, C. 미사용 변수 _ 처리, D. config.* 코드젠)
- Phase 025 ✅ (pkg/queue 패턴 참조)

## SSaC 코드젠 현황 (022 완료 기준)

SSaC가 생성하는 코드:

**@auth:**
```go
if _, err := authz.Check(authz.CheckRequest{Action: "AcceptProposal", Resource: "gig", UserID: currentUser.ID, ResourceID: gig.ClientID}); err != nil {
    c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
    return
}
```

- `authz.Check()` 패키지 함수 호출 (기존 `h.Authz.Check()` 메서드 호출이 아님)
- import `"authz"` 생성

**config.*:**
```go
config.Get("SMTP_HOST")
```

- import `"config"` 생성

## 1. `pkg/authz` 구현

SSaC @auth가 `authz.Check(authz.CheckRequest{...})` 패키지 함수를 호출하므로 **패키지 함수 패턴** 사용.

### API

```go
package authz

import (
    "context"
    "database/sql"
    _ "embed"
    "fmt"
    "os"

    "github.com/open-policy-agent/opa/v1/rego"
)

//go:embed authz.rego
var policyRego string

type CheckRequest struct {
    Action     string
    Resource   string
    UserID     int64
    ResourceID int64
}

type CheckResponse struct{}

var globalEval *rego.PreparedEvalQuery
var globalDB *sql.DB

// Init initializes the global authz evaluator.
func Init(db *sql.DB) error {
    query, err := rego.New(
        rego.Query("data.authz.allow"),
        rego.Module("policy.rego", policyRego),
    ).PrepareForEval(context.Background())
    if err != nil {
        return fmt.Errorf("OPA init failed: %w", err)
    }
    globalEval = &query
    globalDB = db
    return nil
}

// Check evaluates the OPA policy. Returns error if denied.
func Check(req CheckRequest) (CheckResponse, error) {
    if os.Getenv("DISABLE_AUTHZ") == "1" {
        return CheckResponse{}, nil
    }
    if globalEval == nil {
        return CheckResponse{}, fmt.Errorf("authz not initialized")
    }

    opaInput := map[string]interface{}{
        "user":        map[string]interface{}{"id": req.UserID},
        "action":      req.Action,
        "resource":    req.Resource,
        "resource_id": req.ResourceID,
    }

    results, err := globalEval.Eval(context.Background(), rego.EvalInput(opaInput))
    if err != nil {
        return CheckResponse{}, fmt.Errorf("OPA eval failed: %w", err)
    }
    if len(results) == 0 {
        return CheckResponse{}, fmt.Errorf("forbidden")
    }
    allowed, ok := results[0].Expressions[0].Value.(bool)
    if !ok || !allowed {
        return CheckResponse{}, fmt.Errorf("forbidden")
    }
    return CheckResponse{}, nil
}
```

- ownership lookup은 제거 — authz.rego 정책 자체에서 처리
- Init(db) 호출은 main.go에서 생성 (기존 `authz.New(conn)` 패턴과 유사)
- SSaC @auth는 funcspec 없이도 동작 — pkg/authz의 CheckRequest가 interface 역할

### 파일

| 파일 | 내용 |
|---|---|
| `pkg/authz/authz.go` | 싱글턴 API + Check 함수 |
| `pkg/authz/authz_test.go` | 단위 테스트 (DISABLE_AUTHZ=1 우회 등) |

## 2. `pkg/config` 구현

```go
package config

import "os"

// Get returns the environment variable value for the given key.
func Get(key string) string {
    return os.Getenv(key)
}

// MustGet returns the environment variable value, panics if empty.
func MustGet(key string) string {
    v := os.Getenv(key)
    if v == "" {
        panic("required env var not set: " + key)
    }
    return v
}
```

### 파일

| 파일 | 내용 |
|---|---|
| `pkg/config/config.go` | Get + MustGet |
| `pkg/config/config_test.go` | 단위 테스트 |

## 3. `fullend.yaml` 파싱

`ProjectConfig`에 `Authz` 필드 추가:

```go
type ProjectConfig struct {
    // ... 기존 필드
    Authz *AuthzConfig `yaml:"authz"`
}

type AuthzConfig struct {
    Package string `yaml:"package"` // 사용자 커스텀 authz 패키지 경로
}
```

fullend.yaml 예시:

```yaml
# 기본값: github.com/geul-org/fullend/pkg/authz
authz:
  package: github.com/gigbridge/api/internal/authz   # 커스텀 패키지
```

Fallback:
1. `authz.package` 지정 → 해당 경로
2. 미지정 → `github.com/geul-org/fullend/pkg/authz`

### 파일

| 파일 | 변경 |
|---|---|
| `internal/projectconfig/projectconfig.go` | `AuthzConfig` 타입 + `Authz` 필드 추가 |

## 4. authzgen 전환

### 현재 (`internal/gluegen/authzgen.go`)

OPA Authorizer를 고정 `Input{ID interface{}}` + `Check(user, action, resource, input)` 메서드로 자동생성.

### 변경 후

authzgen의 역할 변경:
1. `.rego` 파일 복사 (기존과 동일)
2. `pkg/authz`를 artifacts의 authz 디렉토리에 복사 (커스텀 패키지가 아닌 경우)
3. 커스텀 패키지 사용 시 → .rego 파일만 복사, Go 코드는 생성하지 않음

```go
func GenerateAuthzPackage(policies []*policy.Policy, artifactsDir, modulePath string, authzPackage string) error {
    authzDir := filepath.Join(artifactsDir, "backend", "internal", "authz")
    os.MkdirAll(authzDir, 0755)

    // 1. .rego 파일 복사 (항상)
    for _, p := range policies {
        data, _ := os.ReadFile(p.File)
        os.WriteFile(filepath.Join(authzDir, filepath.Base(p.File)), data, 0644)
    }

    // 2. 커스텀 패키지가 아니면 pkg/authz 기반 코드 생성
    if authzPackage == "" {
        src := generateDefaultAuthzSource(policies)
        os.WriteFile(filepath.Join(authzDir, "authz.go"), []byte(src), 0644)
    }
    // 커스텀 패키지면: 사용자 코드가 이미 존재한다고 가정, .rego만 복사
    return nil
}
```

### 파일

| 파일 | 변경 |
|---|---|
| `internal/gluegen/authzgen.go` | `authzPackage` 파라미터 추가 + 코드 생성 로직 변경 |

## 5. gluegen 변경

### 5-1. authz 참조 변환 제거

SSaC 022가 `authz.Check(authz.CheckRequest{...})` 패키지 함수 호출을 생성하므로:
- `transformSource()`에서 `authz.Check(currentUser,` → `h.Authz.Check(currentUser,` 치환 **제거**
- `"authz"` import → `"{modulePath}/internal/authz"` 변환은 **유지** (authz 패키지 경로 수정)

### 5-2. config import 변환 추가

SSaC가 `"config"` import를 생성하므로, fullend `pkg/config` 경로로 변환:

```go
// Fix config import: "config" → "github.com/geul-org/fullend/pkg/config"
if strings.Contains(src, "\t\"config\"\n") {
    src = strings.ReplaceAll(src, "\t\"config\"\n", "\t\"github.com/geul-org/fullend/pkg/config\"\n")
}
```

기존 `"queue"` → `pkg/queue` 변환과 동일한 패턴.

### 5-3. Handler/Server struct에서 Authz 필드 제거

SSaC가 `h.Authz.Check()` 대신 `authz.Check()` 패키지 함수를 호출하므로:
- `domain.go`: Handler struct에서 `Authz model.Authorizer` 필드 제거
- `server.go`: Server struct에서 `Authz Authorizer` 필드 제거
- main.go 생성: `authz.New(conn)` → `authz.Init(conn)` 변경 (싱글턴 초기화)

### 5-4. GlueInput 확장

```go
type GlueInput struct {
    // ... 기존 필드
    AuthzPackage string // from fullend.yaml, "" = default pkg/authz
}
```

### 파일

| 파일 | 변경 |
|---|---|
| `internal/gluegen/gluegen.go` | authz 치환 제거 + config import 변환 추가 + `AuthzPackage` 필드 |
| `internal/gluegen/domain.go` | Handler Authz 필드 제거 + main.go authz.Init() |
| `internal/gluegen/server.go` | Server Authz 필드 제거 |
| `internal/gluegen/main_go.go` | authz.Init() 호출 생성 |

## 6. crosscheck 변경

### @auth inputs → authz CheckRequest 필드 매칭

`@auth` inputs의 키가 authz CheckRequest struct 필드에 존재하는지 검증.
기본 CheckRequest 필드: `Action`, `Resource`, `UserID`, `ResourceID`.

| Rule | Level | 설명 |
|---|---|---|
| `@auth` inputs 키 → CheckRequest 필드 존재 | ERROR | 필드명 불일치 |
| fullend.yaml `authz` 미설정 + `@auth` 사용 | — | 기본 `pkg/authz` 사용 (에러 아님) |

구현:

```go
func CheckAuthz(funcs []ssacparser.ServiceFunc, authzPackage string, fullendPkgSpecs []funcspec.FuncSpec) []CrossError {
    // authz 패키지의 CheckRequest 필드를 fullendPkgSpecs 또는 기본값에서 추출
    checkRequestFields := map[string]bool{
        "Action": true, "Resource": true,
        "UserID": true, "ResourceID": true,
    }

    for _, fn := range funcs {
        for _, seq := range fn.Sequences {
            if seq.Type != "auth" { continue }
            for key := range seq.Inputs {
                if !checkRequestFields[key] {
                    // ERROR: 필드 불일치
                }
            }
        }
    }
}
```

### 파일

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/authz.go` | 신규 — @auth ↔ authz CheckRequest 교차 검증 |
| `internal/crosscheck/authz_test.go` | 신규 — 테스트 |
| `internal/crosscheck/crosscheck.go` | `AuthzPackage` 필드 + `CheckAuthz` 호출 추가 |

## 7. orchestrator 변경

`CrossValidateInput`과 `GlueInput`에 `AuthzPackage` 전달.

genGlue():
```go
var authzPackage string
if cfg.Authz != nil {
    authzPackage = cfg.Authz.Package
}
input.AuthzPackage = authzPackage
```

genAuthz():
```go
func genAuthz(policyDir, specsDir, artifactsDir string) reporter.StepResult {
    // ... 기존 코드
    var authzPackage string
    if cfg, err := projectconfig.Load(specsDir); err == nil && cfg.Authz != nil {
        authzPackage = cfg.Authz.Package
    }
    gluegen.GenerateAuthzPackage(policies, artifactsDir, modulePath, authzPackage)
}
```

### 파일

| 파일 | 변경 |
|---|---|
| `internal/orchestrator/gen.go` | `genAuthz()` + `genGlue()` 에 authzPackage 전달 |
| `internal/orchestrator/validate.go` | `runCrossValidate()` 에 authzPackage 전달 |

## 전체 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `pkg/authz/authz.go` | 신규 — 기본 OPA 기반 인가 (싱글턴) |
| `pkg/authz/authz_test.go` | 신규 — 테스트 |
| `pkg/config/config.go` | 신규 — 환경변수 조회 |
| `pkg/config/config_test.go` | 신규 — 테스트 |
| `internal/projectconfig/projectconfig.go` | `AuthzConfig` 추가 |
| `internal/gluegen/authzgen.go` | authzPackage 파라미터 + 코드 생성 전환 |
| `internal/gluegen/gluegen.go` | authz 치환 제거 + config import 변환 + AuthzPackage 필드 |
| `internal/gluegen/domain.go` | Handler Authz 필드 제거 + main.go authz.Init() |
| `internal/gluegen/server.go` | Server Authz 필드 제거 |
| `internal/gluegen/main_go.go` | authz.Init() 호출 생성 |
| `internal/crosscheck/crosscheck.go` | AuthzPackage 필드 + CheckAuthz 호출 |
| `internal/crosscheck/authz.go` | 신규 — @auth ↔ CheckRequest 교차 검증 |
| `internal/crosscheck/authz_test.go` | 신규 — 테스트 |
| `internal/orchestrator/gen.go` | authzPackage 전달 |
| `internal/orchestrator/validate.go` | authzPackage 전달 |

## 검증 방법

```bash
# 1. pkg/authz 단위 테스트
go test ./pkg/authz/...

# 2. pkg/config 단위 테스트
go test ./pkg/config/...

# 3. crosscheck 단위 테스트
go test ./internal/crosscheck/...

# 4. 전체 빌드
go build ./cmd/fullend/

# 5. dummy-gigbridge end-to-end
fullend validate specs/gigbridge
fullend gen specs/gigbridge artifacts/gigbridge
cd artifacts/gigbridge/backend && go build ./...
```
