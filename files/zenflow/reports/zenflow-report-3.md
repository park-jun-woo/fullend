# ZenFlow 더미 프로젝트 리포트 #3

## 개요
- 프로젝트: ZenFlow — 멀티테넌트 워크플로우 자동화 SaaS
- 일시: 2026-03-15
- 시작: 07:44 KST
- 종료: 08:00 KST
- 소요시간: 약 16분

## SSOT 작성 현황

| SSOT | 파일 수 | 상태 |
|------|---------|------|
| fullend.yaml | 1 | 완료 |
| DDL | 5 | 완료 (organizations, users, workflows, actions, executions) |
| sqlc queries | 5 | 완료 |
| OpenAPI | 1 (11 endpoints) | 완료 |
| SSaC | 11 | 완료 |
| Model | 1 | 완료 (package model only) |
| States | 1 (5 transitions) | 완료 |
| Policy | 1 (3 rules) | 완료 |
| STML | 2 pages | 완료 |
| Func Spec | 3 (worker, billing x2) | 완료 |
| Hurl Scenario | 1 scenario + 2 invariant | 완료 |

## Validate 결과
- ERROR: 0
- WARNING: 5 (x-sort 인덱스, claims 미참조, func 역방향 매칭)

## Codegen 결과
모든 코드 산출물 정상 생성:
- sqlc, oapi-gen, ssac-gen, ssac-model, stml-gen, glue-gen, hurl-gen, state-gen, authz-gen, func-gen

## Build 결과
- `go build` 성공 (2회차 — 1회차에서 `@empty int64 nil` 비교 오류 발견 후 SSOT 우회 적용)

## 테스트 결과

| 테스트 | 결과 | 비고 |
|--------|------|------|
| smoke.hurl | PASS (9/9) | DISABLE_AUTHZ=1, DISABLE_STATE_CHECK=1 |
| scenario-happy-path.hurl | PASS (7/7) | 워크플로우 생성→액션 추가→활성화→실행 전체 흐름 |
| invariant-tenant-breach.hurl | SKIP | authz 활성화 시 ownership 매칭 제한 (BUG021) |
| invariant-insufficient-credits.hurl | SKIP | authz 활성화 시 403 반환 (org_id claim 미지원) |

## 발견된 버그

| ID | 단계 | 내용 |
|----|------|------|
| BUG019 | validate | `sqlFileToModel()` 복합 테이블명 inflection 불일치 — `execution_logs` → `Execution_log` (기대: `ExecutionLog`) |
| BUG020 | gen→build | `@empty` on int64 필드에서 `== nil` 비교 코드 생성 (compile error) |
| BUG021 | gen→runtime | BearerAuth 미들웨어에서 커스텀 JWT claims (`org_id` 등) 미매핑 |
| BUG022 | validate | Func → SSaC 역방향 crosscheck에서 참조 매칭 실패 WARNING |

## 우회 방법

| 버그 | 우회 |
|------|------|
| BUG019 | 복합 테이블명(`execution_logs`) → 단일어(`executions`)로 변경 |
| BUG020 | `@empty credits.Balance` → `@call billing.CheckCredits`에 `@error 402` 추가하여 func 내부에서 검증 |
| BUG021 | claims에서 `OrgID` 제거, `@get User me = User.FindByID({ID: currentUser.ID})` → `me.OrgID` 사용 |

## 소감

fullend의 SSOT → codegen → test 파이프라인은 전반적으로 잘 동작한다. 특히:
- validate의 교차 검증이 매우 세밀하여 SSOT 간 불일치를 빠르게 잡아준다
- 코드 생성 품질이 높아 별도 수정 없이 빌드 가능 (버그 우회 후)
- smoke.hurl 자동 생성이 편리하다

개선 필요 사항:
1. 커스텀 JWT claims 지원 (가장 큰 제약)
2. `@empty` 타입 인식 개선
3. 복합 테이블명 inflection 통일
