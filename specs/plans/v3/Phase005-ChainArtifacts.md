# Phase 005: Feature Chain에 Artifacts 노드 통합 ✅ 완료

## 목표

`fullend chain` 출력에 SSOT 노드뿐 아니라 파생된 artifacts 함수와 소유권 상태를 포함한다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/orchestrator/chain.go` | 수정 — artifacts 스캔 + chain 결과에 artifact 노드 추가 |
| `internal/reporter/chain.go` | 수정 — Artifacts 섹션 출력 |

## 상세 설계

### chain 흐름 확장

기존 chain: `operationId → OpenAPI, SSaC, DDL, Rego, States, Func, Gherkin, STML`

확장: 기존 chain 실행 후, artifacts 디렉토리에서 해당 operationId에 연결된 함수를 탐색.

```go
// chain.go
func Chain(specsDir, operationID string) ([]ChainLink, error) {
    // 기존 SSOT 탐색 (변경 없음)
    links := traceSSOTs(specsDir, operationID)

    // artifacts 탐색 (신규)
    artifactsDir := inferArtifactsDir(specsDir)  // specs/../artifacts/<project>
    if artifactsDir != "" {
        artifactLinks := traceArtifacts(artifactsDir, operationID)
        links = append(links, artifactLinks...)
    }
    return links, nil
}
```

### artifact 노드 탐색

artifacts 디렉토리의 Go 파일에서 `//fullend:` 디렉티브를 스캔. operationId와 매칭되는 SSOT 경로를 가진 함수를 추출.

매칭 규칙:
- Handler: `ssot=` 경로의 SSaC 파일이 해당 operationId의 것인지 확인
- Model: SSaC가 참조하는 DDL 테이블의 모델 구현 함수
- State: SSaC가 참조하는 stateDiagram의 생성 코드

### 출력

```
── Feature Chain: CreateGig ──

  OpenAPI    api/openapi.yaml:45         POST /gigs
  SSaC       service/gig/create_gig.ssac @post @response
  DDL        db/gigs.sql:1               CREATE TABLE gigs
  Rego       policy/authz.rego:12        resource: gig
  Gherkin    scenario/gig.feature:5      Scenario: Create a new gig

  ── Artifacts ──
  Handler    internal/service/gig/create_gig.go:CreateGig    preserve ✎
  Model      internal/model/gig.go:Create                    preserve ✎
  Authz      internal/authz/authz.go                         gen
  Types      internal/model/types.go:Gig                     gen
```

ChainLink에 `Ownership` 필드 추가:

```go
type ChainLink struct {
    Kind      string // "OpenAPI", "SSaC", ..., "Handler", "Model", ...
    File      string
    Line      int
    Summary   string
    Ownership string // "", "gen", "preserve" (SSOT 노드는 빈 문자열)
}
```

## 의존성

- Phase 001 (`internal/contract` — Directive 파싱)
- Phase 004 (`internal/contract/scan.go` — artifacts 스캔)

## 검증

1. `fullend chain CreateGig specs/` — Artifacts 섹션 출력
2. artifacts 없을 때 — Artifacts 섹션 스킵
3. preserve 함수 — `preserve ✎` 표시
4. gen 함수 — `gen` 표시
5. `go test ./...` 통과
