# ✅ Phase 020: crosscheck 정밀도 개선

## 설계 결정

### SSaC/STML/fullend 책임 경계

- **SSaC** = 서비스 흐름 DSL 파서 + 문법 검증 + 핸들러 코드젠. 코드젠에 OpenAPI를 참조하므로 SSaC validate에 OpenAPI 포함. OPA/stateDiagram은 SSaC 관심사 아님.
- **STML** = UI DSL 파서 + 코드젠.
- **fullend** = 나머지 전부. 특히 Policy ↔ States 교차 검증은 fullend 전담.

## 목표

`fullend validate` crosscheck 정밀도 개선 (11건 false-positive 제거):

| # | 카테고리 | 건수 | 문제 |
|---|---|---|---|
| 1 | @result ↔ DDL | 4 | 외부 패키지 DTO 타입을 DDL 테이블로 찾으려는 오류 |
| 2 | Scenario ↔ States | 4 | 이벤트가 여러 다이어그램에 걸칠 때 상태 추적 실패 |
| 3 | Policy ↔ States | 3 | 다이어그램 ID로 resource를 추론하는 설계 오류 |

## 문제 분석 및 수정

### 1. @result ↔ DDL (4건)

`@call` 시퀀스의 `@result` 타입 (`auth.IssueTokenResponse`, `billing.HoldEscrowResponse` 등)이 DDL 테이블 매칭 대상으로 잡힘. `@call` = 순수 로직 = I/O 금지 = DB 무관이므로 `@call`의 결과 타입은 원칙적으로 DDL과 무관.

**수정**: `CheckSSaCDDL()`에서 `seq.Type == "call"`이면 @result ↔ DDL 체크 자체를 스킵.

```go
// ssac_ddl.go — 수정
for i, seq := range fn.Sequences {
    if seq.Type == "call" {
        continue  // @call = 순수 로직, DDL 무관
    }
    if seq.Result != nil && seq.Result.Type != "" {
        errs = append(errs, checkResultType(seq, st, ctx, i, dtoTypes)...)
    }
}
```

### 2. Scenario ↔ States (4건)

`fullend validate`가 Gherkin 시나리오와 stateDiagram을 교차 검증할 때, 시나리오 스텝들이 상태 전이 순서를 올바르게 따르는지 시뮬레이션함.

`AcceptProposal`이 gig(`open→in_progress`)과 proposal(`pending→accepted`) 두 다이어그램에 존재. `eventDiagram`이 `map[string]*StateDiagram` (1:1)이라 마지막으로 등록된 다이어그램이 이전 것을 덮어씀. gig 다이어그램이 무시되어 AcceptProposal 후에도 gig 상태가 `open`에서 전진 못함. 이후 SubmitWork(`in_progress`에서만 가능), ApproveWork(`under_review`에서만 가능) 스텝에서 false-positive WARN 발생.

**수정**: `eventDiagram`을 `map[string][]*StateDiagram` (1:N)으로 변경. 이벤트 발생 시 모든 관련 다이어그램의 상태를 동시에 전진.

```go
// scenario.go — 현재 (1:1, 덮어쓰기 문제)
eventDiagram := make(map[string]*statemachine.StateDiagram)
for _, d := range diagrams {
    for _, ev := range d.Events() {
        eventDiagram[ev] = d  // 같은 키면 덮어씀
    }
}

// 수정 (1:N, 모든 다이어그램 유지)
eventDiagrams := make(map[string][]*statemachine.StateDiagram)
for _, d := range diagrams {
    for _, ev := range d.Events() {
        eventDiagrams[ev] = append(eventDiagrams[ev], d)
    }
}

// 검증/전진 시 모든 관련 다이어그램에 대해 수행
for _, d := range eventDiagrams[step.OperationID] {
    state := currentState[d.ID]
    // 유효성 체크 + 상태 전진
}
```

### 3. Policy ↔ States (3건)

다이어그램 ID(`d.ID`)로 resource를 추론하지만, 다이어그램 ID ≠ SSaC authorize resource:
- `submit_work` → Rego resource `gig_assignee`, 다이어그램 `gig` → 잘못된 추론 `("submit_work", "gig")`
- `accept` → Rego resource `gig`, 다이어그램 `proposal` → 잘못된 추론 `("accept", "proposal")`
- `reject` → Rego resource `gig`, 다이어그램 `proposal` → 잘못된 추론 `("reject", "proposal")`

**수정**: Policy ↔ States 섹션 자체를 제거. 이 검증은 이미 두 경로로 커버됨:
- **States ↔ SSaC**: 다이어그램 이벤트에 대응하는 SSaC 함수 존재 확인
- **Policy ↔ SSaC**: SSaC @auth (action, resource) 쌍이 Rego에 존재 확인

전이적으로 `다이어그램 이벤트 → SSaC 함수 → Rego 규칙` 검증이 완성되므로 Policy ↔ States는 중복이며, d.ID 기반 추론으로 인해 버그만 발생.

```go
// policy.go — 삭제 대상 (136-168행)
// --- Policy ↔ States --- 섹션 전체 제거
```

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/ssac_ddl.go` | `CheckSSaCDDL()`: `seq.Type == "call"` 시 @result ↔ DDL 체크 스킵 |
| `internal/crosscheck/scenario.go` | `checkScenarioStates()`: eventDiagram 1:N 변경 |
| `internal/crosscheck/policy.go` | Policy ↔ States 섹션 (136-168행) 제거 |

## 의존성

없음 (fullend 자체 crosscheck 로직만 수정)

## 검증 방법

```bash
go test ./internal/crosscheck/...
fullend validate specs/dummy-gigbridge  # false-positive 11건 제거 확인
```
