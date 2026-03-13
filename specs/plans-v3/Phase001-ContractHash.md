# Phase 001: Contract Hash 계산 + 디렉티브 인프라 ✅ 완료

## 목표

SSOT에서 계약 해시를 계산하고, `//fullend:` 디렉티브를 파싱/생성하는 기반 패키지를 만든다.

## 배경

Contract-Based Code Generation의 핵심 전제: 계약은 Go 시그니처가 아니라 SSOT 명세에서 파생된다. 이 Phase에서 해시 계산 로직과 디렉티브 구조체를 먼저 구현한다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/contract/directive.go` | 신규 — Directive 구조체, 파싱, 직렬화 |
| `internal/contract/hash.go` | 신규 — SSOT별 contract hash 계산 |
| `internal/contract/directive_test.go` | 신규 — 디렉티브 파싱/직렬화 테스트 |
| `internal/contract/hash_test.go` | 신규 — hash 계산 테스트 |

## 상세 설계

### Directive 구조체

```go
package contract

type Directive struct {
    Ownership string // "gen" or "preserve"
    SSOT      string // SSOT 파일 상대경로
    Contract  string // 7자리 SHA256
}

func Parse(comment string) (*Directive, error)    // "//fullend:gen ssot=... contract=..." → Directive
func (d *Directive) String() string               // Directive → "//fullend:gen ssot=... contract=..."
```

### Contract Hash 계산

각 SSOT 종류별로 해싱 대상이 다름:

```go
// Service Handler: operationId + 시퀀스 타입 + request fields + response fields
func HashServiceFunc(sf ssacparser.ServiceFunc) string

// Model Implementation: 함수명 + 파라미터 타입 + 반환 타입
func HashModelMethod(name string, params []string, returns []string) string

// State Machine: state 목록 + transition 목록
func HashStateDiagram(sd *statemachine.StateDiagram) string

// Middleware: CurrentUser struct fields
func HashClaims(claims map[string]string) string
```

모두 SHA256 → 앞 7자리 hex 반환.

## 의존성

- `internal/statemachine` — StateDiagram 구조체
- `ssac/parser` — ServiceFunc 구조체

## 검증

```bash
go test ./internal/contract/...
```

- Parse("//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1") → 정확한 Directive
- Parse("//fullend:preserve ...") → Ownership == "preserve"
- 잘못된 형식 → 에러
- HashServiceFunc 동일 입력 → 동일 해시
- HashServiceFunc 다른 입력 → 다른 해시
