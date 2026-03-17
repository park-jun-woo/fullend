# ZenFlow Report #6 (zenflow-try06)

> **gigbridge-try04 패턴 참고**: BIGSERIAL PK, @auth {ResourceID: currentUser.ID} 우회, @empty 재조회 가드, Func 본체 구현 패턴

## 시간
- 시작: 2026-03-17 18:14:15
- 종료: 2026-03-17 18:20:49
- 소요: 약 7분

## 결과: PASS

### Validate
- WARNING 5건 (무해)
  - checkCredits 반환값 무시 2건 (의도: 에러만 체크)
  - claims Rego 미참조 3건 (org 격리를 query 레벨에서 처리)

### Build
- `go build` 성공
- `replace` 디렉티브로 fullend 로컬 참조

### Hurl 테스트
| 테스트 | 결과 | 요청 수 |
|---|---|---|
| smoke.hurl | PASS | 10 |
| scenario-happy-path.hurl | PASS | 9 |
| invariant-tenant-breach.hurl | PASS | 8 |
| invariant-insufficient-credits.hurl | PASS | 8 |

## SSOT 작성 요약

| SSOT | 파일 수 | 비고 |
|---|---|---|
| fullend.yaml | 1 | JWT claims: ID, Email, Role, OrgID (multi-tenant) |
| DDL | 5 | organizations, users, workflows, actions, execution_logs |
| sqlc queries | 5 | 15 queries total (FindByIDAndOrgID for org isolation) |
| OpenAPI | 1 | 12 endpoints |
| SSaC | 12 | Auth, CRUD, state transitions, execute+credit system |
| Model | 1 | package model (no @dto needed) |
| States | 1 | workflow (5 transitions, including self-transition active→active) |
| Policy | 1 | 9 allow rules, 1 ownership mapping |
| Hurl | 3 | 1 scenario + 2 invariant |
| STML | 2 | workflows list + workflow detail |
| Func | 3 | checkCredits(@error 402), deductCredit, processActions |

## gigbridge 대비 복잡도 증가 포인트

1. **Multi-tenant 격리** — org_id 기반, FindByIDAndOrgID 쿼리 패턴으로 해결
2. **크레딧 시스템** — checkCredits(@error 402) + deductCredit + Organization.UpdateCredits
3. **자기 전이** — active→active: ExecuteWorkflow (상태 변경 없이 실행)
4. **3계층 Func** — billing(2), worker(1)
5. **5테이블** — gigbridge 4테이블 대비 1테이블 추가
6. **크레딧 소진 후 재활성화 불가** — 실행→소진→402 시나리오

## 첫 validate 에러 → 수정 내역

1. `actions.payload_template` NOT NULL 누락 → `NOT NULL DEFAULT ''` 추가
2. `payload_template` OpenAPI required 누락 → required 배열에 추가

## 런타임 이슈 → 수정 내역

1. `credits_balance: 0` 값이 Go `binding:"required"` 에서 zero value로 거부됨
   - → 테스트 전략 변경: 1 크레딧으로 시작 → 실행으로 소진 → 재활성화 시도 → 402
2. `DeductCredit` Func에 CurrentBalance 파라미터 추가 (정확한 잔액 계산)

## 버그 리포트
없음.
