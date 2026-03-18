# Phase003 — pkg/authz 런타임 교체: OPA → toulmin Graph Builder

## 목표

`pkg/authz/`에서 OPA Rego 런타임 평가를 제거하고 `toulmin.Graph.Evaluate()`로 교체한다. 외부 의존성 `github.com/open-policy-agent/opa` 제거.

## 배경

현재 `pkg/authz/check.go`는 런타임에 Rego 텍스트를 OPA에 로딩하여 `data.authz.allow`를 평가한다. toulmin Graph Builder의 `Evaluate(claim, ground)` 로 대체한다.

## 현재 구조

```
pkg/authz/
├── init_authz.go        # OPA_POLICY_PATH에서 Rego 로딩, globalPolicy 저장
├── check.go             # rego.New() → Eval() → allow 판정
├── check_request.go     # CheckRequest 구조체
├── check_response.go    # CheckResponse 구조체
├── load_owners.go       # DB에서 ownership 데이터 로딩
├── ownership_mapping.go # OwnershipMapping 구조체
```

## 새 구조

```
pkg/authz/
├── init_authz.go        # toulmin Graph 초기화 (생성된 코드에서 Graph 주입)
├── check.go             # graph.Evaluate(claim, ground) → verdict 판정
├── check_request.go     # CheckRequest 구조체 (유지)
├── check_response.go    # CheckResponse 구조체 (유지)
├── load_owners.go       # DB에서 ownership 데이터 로딩 (유지)
├── ownership_mapping.go # OwnershipMapping 구조체 (유지)
```

## 변경 상세

### Init 변경

```go
// 기존
var globalPolicy string // Rego 텍스트

func Init(db *sql.DB, mappings []OwnershipMapping) error {
    raw, _ := os.ReadFile(os.Getenv("OPA_POLICY_PATH"))
    globalPolicy = string(raw)
    // ...
}

// 신규
var globalGraph *toulmin.GraphBuilder

func Init(db *sql.DB, mappings []OwnershipMapping, graph *toulmin.GraphBuilder) error {
    globalGraph = graph
    // ... ownership/DB 초기화는 동일
}
```

### Check 변경

```go
// 기존
func Check(req CheckRequest) (CheckResponse, error) {
    query, _ := rego.New(
        rego.Query("data.authz.allow"),
        rego.Module("policy.rego", globalPolicy),
        rego.Store(store),
        rego.Input(opaInput),
    ).Eval(context.Background())
    allowed := query[0].Expressions[0].Value.(bool)
    // ...
}

// 신규
func Check(req CheckRequest) (CheckResponse, error) {
    if os.Getenv("DISABLE_AUTHZ") == "1" {
        return CheckResponse{}, nil
    }

    results, err := globalGraph.Evaluate(req, nil)
    if err != nil {
        return CheckResponse{}, fmt.Errorf("authz evaluation failed: %w", err)
    }

    // DenyByDefault가 warrant. verdict > 0이면 deny 유지(허용 규칙이 무력화 실패)
    for _, r := range results {
        if r.Verdict > 0 {
            return CheckResponse{}, fmt.Errorf("denied by %s (verdict: %.2f)", r.Name, r.Verdict)
        }
    }
    return CheckResponse{}, nil
}
```

### CheckTrace 추가 (디버깅용)

```go
func CheckTrace(req CheckRequest) (CheckResponse, []toulmin.EvalResult, error) {
    results, err := globalGraph.EvaluateTrace(req, nil)
    if err != nil {
        return CheckResponse{}, nil, err
    }
    // verdict 판정 + trace 반환
    for _, r := range results {
        if r.Verdict > 0 {
            return CheckResponse{}, results, fmt.Errorf("denied by %s", r.Name)
        }
    }
    return CheckResponse{}, results, nil
}
```

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `pkg/authz/init_authz.go` | OPA 로딩 → toulmin Graph 주입 |
| `pkg/authz/check.go` | rego.New().Eval() → graph.Evaluate() |
| `pkg/authz/check_trace.go` (신규) | EvaluateTrace 기반 디버깅 함수 |
| `go.mod` | `github.com/open-policy-agent/opa` 제거, `github.com/park-jun-woo/toulmin` 추가 |

## 의존성

- `github.com/park-jun-woo/toulmin/pkg/toulmin` — GraphBuilder, EvalResult, TraceEntry

## 검증 방법

- `go test ./pkg/authz/...` 통과
- OPA 관련 import가 전체 프로젝트에서 0건인지 확인: `grep -r "open-policy-agent" .`
- `go mod tidy` 후 OPA 의존성 제거 확인
