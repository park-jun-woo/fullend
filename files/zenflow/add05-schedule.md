# ZenFlow Add-on #05 — 스케줄 실행 (session + @publish 연계)

## 개요
워크플로우에 cron 스케줄을 등록하면 지정 시간에 자동 실행. session에 다음 실행 시각을 저장하고, @publish로 실행 이벤트 발행.

## 검증 포인트 (미검증 fullend 기능)
- **pkg/session**: `session.Session.Set/Get/Delete` — 패키지 프리픽스 @model
- **session.backend: postgres** 설정
- **@publish + 기존 ExecuteWorkflow 재사용**: 스케줄 등록 → 실행 시점에 @publish "workflow.schedule.trigger" 발행

## fullend.yaml 변경
- `session.backend: postgres` 추가

## 신규 엔드포인트
- **POST /workflows/{id}/schedule** (`SetSchedule`): cron 표현식으로 스케줄 등록, session에 다음 실행 시각 저장
- **GET /workflows/{id}/schedule** (`GetSchedule`): 현재 스케줄 조회 (session에서 읽기)
- **DELETE /workflows/{id}/schedule** (`DeleteSchedule`): 스케줄 해제 (session에서 삭제)

## DDL
- 없음 — session 모델이 스토리지 담당 (fullend_sessions 테이블 자동)

## SSaC 설계
- SetSchedule: `@get Workflow` → org 격리 → `@call schedule.ParseCron({Expression: request.cron})` → `@post session.Session.Set({key: scheduleKey, value: cronExpr, ttl: 0})` → `@response`
- GetSchedule: `@get session.Session.Get({key: scheduleKey})` → `@response`
- DeleteSchedule: `@put session.Session.Delete({key: scheduleKey})` → `@response`

## Custom Functions
- `schedule.ParseCron(Expression)`: cron 표현식 유효성 검사 + 다음 실행 시각 계산 (purity 준수)

## E2E Scenario
- 워크플로우 생성/활성화 → 스케줄 등록 → 스케줄 조회 확인 → 스케줄 해제 → 조회 시 빈 값 확인
