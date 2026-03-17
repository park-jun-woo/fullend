# ZenFlow Add-on #04 — 실행 결과 파일 첨부

## 개요
ExecuteWorkflow 결과를 파일로 저장. pkg/file의 FileModel 인터페이스 검증.

## fullend.yaml 변경
- `file.backend: local` 추가

## 신규 엔드포인트
- **POST /workflows/{id}/execute-with-report** (`ExecuteWithReport`): 실행 + 결과 파일 생성/업로드
- **GET /execution-logs/{id}/report** (`GetExecutionReport`): 실행 로그에 연결된 파일 다운로드

## DDL 변경
- `execution_logs` 테이블에 `report_key VARCHAR(255) NOT NULL DEFAULT ''` 추가 — 파일 스토리지 키

## SSaC 설계
- ExecuteWithReport: 기존 ExecuteWorkflow 흐름 + `@call report.GenerateReport({...})` → `@post file.File.Upload({key: reportKey, body: ...})` → execution_log에 report_key 저장
- GetExecutionReport: `@get ExecutionLog` → `@get file.File.Download({key: log.ReportKey})`

## Custom Functions
- `report.GenerateReport(WorkflowID, ActionCount, Status)`: 실행 결과를 텍스트로 포매팅 (purity 준수)

## 검증 포인트
- **pkg/file 패키지 모델**: `file.File.Upload`, `file.File.Download` — 패키지 프리픽스 @model
- **file.backend: local** 설정
- **SSaC에서 pkg 모델 사용**: `@post file.File.Upload({key: ..., body: ...})` 패턴
- Go interface 파라미터 매칭 (context.Context 생략, key/body 이름 일치)

## E2E Scenario
- 워크플로우 생성/활성화/실행(with report) → 실행 로그에 report_key 확인 → 파일 다운로드 검증
