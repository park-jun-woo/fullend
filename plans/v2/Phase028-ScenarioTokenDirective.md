# Phase 028: Scenario Token Directive

## 목표
Gherkin 시나리오에서 요청별 토큰을 명시적으로 지정할 수 있는 문법 추가.
멀티롤 시나리오에서 토큰 전환 버그 해결.

## 문법
```
When POST AcceptProposal {"id": proposalResult.proposal.id} clientToken -> gigResult
When POST SubmitWork freelancerToken
When POST ApproveWork clientToken
```

토큰 변수명은 `Token` 또는 `token`을 포함해야 함 (`\w*[Tt]oken\w*`).
기존 변수(`gigResult`, `proposal2`)와 구분.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/scenario/types.go` | Step에 Token 필드 추가 |
| `internal/scenario/parser.go` | 정규식에 토큰 그룹 추가 |
| `internal/gluegen/hurl_scenario.go` | 토큰 선택 로직: 명시 > currentToken |
| `internal/crosscheck/scenario.go` | checkTokenRoles에서 명시 토큰 반영 |
| `specs/gigbridge/scenario/gig_lifecycle.feature` | 토큰 지시자 추가 |

## 검증 방법
1. `go test ./internal/scenario/...` — 파서 테스트
2. `go test ./internal/crosscheck/...` — crosscheck 테스트
3. `fullend validate specs/gigbridge` — WARNING 감소 확인
4. `fullend gen` → `hurl --test` — 생성된 Hurl에서 올바른 토큰 사용 확인
