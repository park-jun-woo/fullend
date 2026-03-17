# ZenFlow Try04 Report

## 소요 시간
- 시작: 2026-03-15 11:53:25
- 종료: 2026-03-15 12:06:46
- 소요: 약 13분

## 결과
- fullend validate: PASS (2 warnings — claims Rego 미참조, 의도적)
- fullend gen: PASS (10 endpoints, 10 service files, 2 STML pages)
- go build: PASS (수동 패치 후)
- hurl --test: 4/4 PASS (smoke + scenario + invariant 2)

## SSOT 구성
| SSOT | 파일 수 | 비고 |
|---|---|---|
| fullend.yaml | 1 | go/gin, typescript/react |
| DDL | 5 | organizations, users, workflows, actions, execution_logs |
| sqlc queries | 5 | CRUD queries |
| OpenAPI | 1 | 10 endpoints |
| SSaC | 10 | 2 auth + 8 workflow |
| Model | 1 | package model only |
| States | 1 | workflow (5 transitions) |
| Policy | 1 | 7 Rego allow rules (role-based) |
| Hurl tests | 3 | 1 scenario + 2 invariant |
| STML | 2 | workflows list + detail |
| Func | 2 | billing.validateCredits, worker.processAction |

## 발견된 버그
### BUG024: Claims *_id 필드가 string으로 생성됨
- 매뉴얼: `*_id → int64`이지만 실제로는 항상 `string` 생성
- CurrentUser, IssueTokenRequest, VerifyTokenResponse 모두 영향
- 수동 패치: `string` → `int64` + JWT float64 파싱 추가

### BUG025: @auth 템플릿이 currentUser.ID 하드코딩
- claims 키가 `UserID`여도 생성 코드는 `currentUser.ID` 사용
- `currentUser.UserID`로 생성되어야 함
- 수동 패치: sed로 일괄 치환

## 설계 결정
1. **UUID → BIGSERIAL**: fullend의 authz.CheckRequest가 int64 기반이므로 UUID 대신 BIGSERIAL 사용
2. **Org 격리**: @auth + @ownership 조합이 org_id 기반 격리를 지원하지 않아 role-based @auth + 쿼리 필터링으로 대체
3. **Credits 검증**: Func purity 규칙 때문에 DB 접근 불가 → `validateCredits` func에 잔액을 인자로 전달, 0 이하면 에러(@error 402)
4. **DeductCredit**: 리터럴 정수 1을 SSaC 인자로 전달 시 "변수 미선언" 에러 → sqlc 쿼리에서 `-1` 하드코딩으로 우회
5. **ExecuteWorkflow 루프**: SSaC에 루프 없음 → 액션 목록 조회 후 batch 단위 func 호출로 단순화
