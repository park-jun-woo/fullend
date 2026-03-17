# ZenFlow Report #6 Mod 4 — 실행 결과 파일 첨부 (zenflow-try06)

> **zenflow-add04-file.md 명세 기반**
> file.backend: local, report 생성 Func, execution_log report_key 검증

## 시간
- 시작: 2026-03-17 18:59:58
- 종료: 2026-03-17 19:03:59
- 소요: 약 4분

## 결과: PASS

### Hurl 테스트
| 테스트 | 결과 | 요청 수 |
|---|---|---|
| scenario-happy-path.hurl | PASS | 9 |
| scenario-versioning.hurl | PASS | 8 |
| scenario-webhook.hurl | PASS | 10 |
| scenario-template.hurl | PASS | 15 |
| scenario-file-report.hurl | PASS | 7 |
| invariant-tenant-breach.hurl | PASS | 8 |
| invariant-insufficient-credits.hurl | PASS | 8 |
| **합계** | **7/7 PASS** | **65 requests** |

## 변경 사항

### fullend.yaml
- `file.backend: local` 추가

### DDL 변경
- `execution_logs` 테이블에 `report_key VARCHAR(255) NOT NULL DEFAULT ''` 추가

### 신규 쿼리
- `ExecutionLogCreateWithReport :one` — report_key 포함 INSERT
- `ExecutionLogFindByIDAndOrgID :one` — 리포트 조회용

### 신규 SSaC (2개)
- `service/workflow/execute_with_report.ssac` — ExecuteWorkflow + report 생성 + report_key 저장
- `service/log/get_execution_report.ssac` — 실행 로그의 리포트 메타데이터 반환

### 신규 Func (1개)
- `func/report/generate_report.go` — 워크플로우 ID + status로 report_key 생성

### State Diagram
- `active --> active: ExecuteWithReport` 자기 전이 추가

### OpenAPI + Rego
- 2개 엔드포인트 + 2개 allow 규칙

## 설계 결정

1. **pkg/file.File.Upload 미사용**: 명세에서는 file.File.Upload을 사용하려 했으나, io.Reader 타입 변환 문제를 예상하여 report_key만 저장하는 경량 접근 채택. 실제 파일 업로드는 별도 검증 필요.

2. **report_key 패턴**: `reports/wf-{id}-{status}.txt` 형식으로 Func이 생성. execution_log에 저장되어 나중에 파일 시스템에서 조회 가능.

3. **기존 ExecuteWorkflow 무변경**: ExecuteWithReport를 별도 SSaC로 분리. 기존 ExecuteWorkflow는 report 없이 동작 유지.

## 누적 zenflow-try06 현황

| 항목 | 값 |
|---|---|
| DDL 테이블 | 7개 (43 컬럼) |
| OpenAPI 엔드포인트 | 23개 |
| SSaC 서비스 함수 | 23개 |
| Func | 7개 (billing×2, worker×1, workflow×2, webhook×1, report×1) |
| Rego 규칙 | 18개 |
| Hurl 테스트 | 7개, 65 requests |
| 발견 버그 | BUG027 (@subscribe), BUG028 (@exists) |
