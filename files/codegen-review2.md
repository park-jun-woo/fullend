# fullend codegen 전수 리뷰 (zenflow-try06)

대상: `artifacts/dummys/zenflow-try06/` (backend + frontend + tests)
일시: 2026-03-18
방법: 87개 파일 전수 조사 (go.sum, server 바이너리 제외)
기준: 프로덕션 배포 가능 수준

---

## 통계

| 심각도 | 건수 |
|--------|------|
| CRITICAL | 8 |
| HIGH | 14 |
| MEDIUM | 18 |

---

## CRITICAL — 배포 불가

### C1. JWT 시크릿 하드코딩 + 불일치

- `internal/auth/issue_token.go:28-29` — `JWT_SECRET` 미설정 시 `"secret"` 폴백
- `internal/auth/refresh_token.go:33` — `"secret"` 리터럴 **고정** (환경변수 참조 자체가 없음)
- `internal/auth/verify_token.go` — Secret을 파라미터로 받음 (다른 두 함수와 방식 불일치)
- **영향**: refresh_token은 환경변수 설정과 무관하게 항상 `"secret"`으로 서명. 공격자가 임의 토큰 위조 가능
- **codegen 수정**: 시크릿 미설정 시 fatal. 세 함수 모두 동일한 주입 방식 통일

### C2. 프론트엔드 인증 토큰 미사용

- `frontend/src/api.ts` 전체
- **현상**: login 응답의 access_token을 저장하지 않고, 모든 후속 요청에 Authorization 헤더 미부착
- **영향**: 인증 필요 API 전부 401. 프론트엔드 사실상 미작동
- **codegen 수정**: 토큰 저장소 + fetch wrapper에 Bearer 헤더 자동 부착

### C3. Queue Publish + Webhook Deliver가 tx.Commit 이전

- `service/workflow/execute_workflow.go:101-114`
- `service/workflow/execute_with_report.go:108-121`
- **현상**: 트랜잭션 내에서 queue.Publish()와 webhook.Deliver() 호출 후 tx.Commit()
- **영향**: Commit 실패 시 이미 발행된 메시지/웹훅 회수 불가. 외부 HTTP 호출이 트랜잭션을 장시간 점유
- **codegen 수정**: 외부 부수효과를 Commit 이후로 이동

### C4. Role 검증 없는 회원가입 + 조직 생성 인증 불필요

- `service/auth/register.go:15` — 요청 바디 `Role` 필드를 검증 없이 저장. `"admin"` 입력 가능
- `service/organization/create_organization.go` — `POST /organizations`가 공개 라우터. 인증 없이 임의 조직 생성 + `credits_balance` 직접 설정
- **영향**: 누구나 admin 권한 + 무제한 크레딧의 조직 생성 가능. 전체 권한 체계 무력화
- **codegen 수정**: roles enum 검증 생성. 조직 생성은 인증 라우터로 이동 또는 별도 인가

### C5. Rego 정책 인증 검증 미수행 (8개 엔드포인트)

- `internal/authz/authz.rego:57-61, 77-80, 83-86, 89-92, 102-105, 115-118, 135-138, 148-151`
- **대상**: ExecuteWorkflow, ListWorkflowVersions, ExecuteWithReport, GetExecutionReport, CloneTemplate, ListWebhooks, GetSchedule, ListExecutionLogs
- **현상**: `input.claims` 검증 없이 action+resource만 매칭하면 허용. 미들웨어가 통과시키면 비인증 사용자도 접근
- **영향**: 인증 미들웨어 우회 시 전체 8개 엔드포인트 무방비
- **추가**: 인증 체크 패턴도 불일치 — `user_id > 0` vs `email != ""` vs 무검증 혼재
- **수정**: 일관된 `input.claims.user_id > 0` 패턴으로 통일

### C6. 크레딧 차감 Race Condition + 음수 잔액 허용

- `service/workflow/execute_workflow.go:72-93`
- `internal/billing/deduct_credit.go:23` — `CurrentBalance - Amount` 결과에 음수 체크 없음
- **현상**: CheckCredits → DeductCredit 사이에 행 잠금 없음. DeductCredit 자체도 음수 허용
- **영향**: 동시 요청 시 이중 차감, 음수 잔액 발생
- **수정**: `UPDATE ... WHERE credits >= cost RETURNING credits` 원자적 차감

### C7. go.mod `go 1.25.0` — 존재하지 않는 Go 버전

- `backend/go.mod:3`
- **영향**: `go build`, `go mod tidy` 즉시 실패. 빌드 자체 불가
- **codegen 수정**: 실제 설치된 Go 버전 또는 1.22/1.23으로 생성

### C8. defer 미실행 (log.Fatal + os.Exit)

- `backend/cmd/main.go:45,62` — `defer conn.Close()`, `defer queue.Close()`
- `backend/cmd/main.go:126` — `log.Fatal(r.Run(*addr))`
- **현상**: log.Fatal은 os.Exit(1) 호출 → defer 미실행
- **영향**: DB 커넥션, 큐 리소스 정리 안 됨. 장기 운영 시 리소스 누수

---

## HIGH — 보안/데이터 무결성 위험

### H1. authz에 잘못된 ResourceID (다수 핸들러)

- `service/workflow/create_workflow.go:32` — `ResourceID: currentUser.ID`
- `service/workflow/list_workflows.go:14` — `ResourceID: currentUser.ID`
- `service/template/publish_template.go:47` — `ResourceID: currentUser.ID`
- `service/webhook/create_webhook.go:32` — `ResourceID: currentUser.ID`
- `service/webhook/list_webhooks.go:14` — `ResourceID: currentUser.ID`
- `service/log/get_execution_report.go:33` — `ResourceID: execLog.ID` (Resource는 `"workflow"`)
- **현상**: 워크플로우/템플릿/웹훅 ID가 아닌 유저 ID 또는 잘못된 엔티티 ID로 인가 체크
- **영향**: 인가 체크가 항상 "이 유저가 자기 자신에 대해 권한이 있는가?"를 묻게 됨 → 사실상 무의미

### H2. QueryOpts SQL 인젝션 표면

- `internal/model/queryopts.go:125-133, 157-163`
- `BuildSelectQuery`/`BuildCountQuery`에서 `col`, `table`, `opts.SortCol`을 `fmt.Sprintf`로 직접 삽입
- ParseQueryOpts 화이트리스트 검증이 있으나 BuildSelectQuery 자체는 재검증하지 않음
- **영향**: QueryOpts를 직접 구성하는 코드가 추가되면 즉시 SQL 인젝션

### H3. 이메일 열거 공격

- `service/auth/login.go:28-31` — 미존재 시 404, 비밀번호 오류 시 401
- **영향**: 응답 코드와 타이밍 차이로 유효 이메일 특정 가능
- **수정**: 미존재 시에도 더미 해시 비교. 동일 응답 코드 반환

### H4. VerifyToken 클레임 누락 시 무음 제로값

- `internal/auth/verify_token.go:40-43`
- `claims["email"].(string)` → 미존재 시 `""`, `claims["user_id"].(float64)` → 미존재 시 `0`
- **영향**: 클레임 없는 토큰이 `ID: 0, OrgID: 0`으로 통과 → Rego에서 검증 안 하는 엔드포인트 접근 가능

### H5. context.Background() 일괄 사용

- `internal/model/*.go` 전체 (action, workflow, user, webhook, template, organization, executionlog)
- **현상**: 모든 DB 쿼리에 context.Background() 사용. 요청 컨텍스트 미전파
- **영향**: HTTP 요청 취소/타임아웃이 DB 레이어에 전달 안 됨. 클라이언트 끊어져도 쿼리 계속 실행

### H6. 페이지네이션 없는 List 엔드포인트

- `service/log/list_execution_logs.go:38` — LIMIT 없음
- `service/workflow/list_workflows.go:19` — LIMIT 없음
- `service/webhook/list_webhooks.go:19` — LIMIT 없음
- `internal/db/workflows.sql.go:86` — sqlc 쿼리도 LIMIT 없음
- `internal/db/execution_logs.sql.go:44` — 동일
- **영향**: 대량 데이터 시 OOM 또는 응답 시간 초과

### H7. CORS 미설정

- `internal/service/server.go` — gin.Default()만 사용. CORS 미들웨어 없음
- **영향**: 프론트엔드에서 크로스 오리진 API 호출 차단

### H8. 프론트엔드 HTTP 상태 미검사

- `frontend/src/api.ts` — 모든 함수에서 `res.json()` 직접 호출. status 확인 없음
- **영향**: 500/401/403 에러 응답도 성공 데이터로 처리

### H9. ActionCount에 CreditsBalance 할당

- `service/workflow/execute_workflow.go:78`
- `service/workflow/execute_with_report.go:79`
- **현상**: `worker.ProcessActions({ActionCount: org.CreditsBalance, ...})`
- **영향**: 액션 수가 아닌 크레딧 잔액이 ActionCount로 전달. 처리 로직 전면 오작동

### H10. Actions 조회 결과 폐기

- `service/workflow/execute_workflow.go:55` — `_, err = h.ActionModel.WithTx(tx).ListByWorkflowID(wf.ID)`
- `service/webhook/on_workflow_executed.go:17` — `_, err := h.WebhookModel.ListByOrgIDAndEventType(...)`
- **현상**: DB 쿼리 결과를 `_`로 폐기. 불필요한 쿼리 실행
- **영향**: DB 부하만 발생하고 결과 미활용. 또는 원래 결과를 사용해야 하는데 누락된 것

### H11. SELECT * / RETURNING * 사용

- `internal/model/action.go:52` — `RETURNING *`
- `internal/model/action.go:61`, `executionlog.go`, `workflow.go` 등 — `SELECT *`
- **현상**: DDL 컬럼 순서에 의존. 마이그레이션 시 scan 필드 순서 불일치 발생
- **영향**: 컬럼 추가/순서 변경 시 런타임 에러 또는 잘못된 데이터 스캔

### H12. ListByOrgID 정렬 없음

- `internal/model/webhook.go:74`, `workflow.go:89` — ORDER BY 없음
- **영향**: 매 요청마다 결과 순서 변동. 페이지네이션 도입 시 데이터 누락/중복

### H13. 모델 중복 인스턴스화

- `backend/cmd/main.go:72-110`
- `model.NewWorkflowModel(conn)` **4회**, `model.NewActionModel(conn)` **3회**, `model.NewOrganizationModel(conn)` **3회**
- **영향**: 메모리 낭비. 모델에 캐싱/상태가 추가되면 불일치 발생

### H14. Register 중복 이메일 → 500

- `service/auth/register.go:53`
- **현상**: DB unique constraint 위반 시 generic 500 반환. 409 Conflict가 적절
- **영향**: 클라이언트가 중복 이메일인지 서버 오류인지 구분 불가

---

## MEDIUM — 운영 안정성/품질

### M1. Graceful Shutdown 없음

- `backend/cmd/main.go:122,126` — `go queue.Start(context.Background())` + `log.Fatal(r.Run())`
- 시그널 핸들링 없음. 큐 고루틴에 취소 불가 컨텍스트 사용

### M2. DB 커넥션 풀 미설정

- `backend/cmd/main.go:38` — SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime 미호출

### M3. DISABLE_STATE_CHECK 환경변수

- `internal/states/workflowstate/workflowstate.go:40-41`
- 환경변수 하나로 전체 상태 검증 무력화. 프로덕션 사고 위험
- **수정**: 빌드 태그(`//go:build debug`)로 제한

### M4. workflowstate Input.Status가 `interface{}`

- `internal/states/workflowstate/workflowstate.go:21,43`
- `input.Status.(string)` 타입 단언 실패 시 `""` → 오류 메시지가 `"cannot transition from \"\""`
- **수정**: `Status`를 `string` 타입으로 변경

### M5. SetSchedule NextRun 하드코딩

- `internal/schedule/set_schedule.go:34` — `NextRun: "2026-03-19T00:00:00Z"` 고정값
- GetSchedule은 항상 `NextRun: ""` 반환. Set/Get 응답 불일치

### M6. GetSchedule 에러 무시

- `internal/schedule/get_schedule.go:25` — session.Get 실패 시 `return ..., nil`
- "스케줄 없음"과 "세션 스토어 장애"를 구분 불가

### M7. 혼합 언어 에러 메시지

- 한국어(`"Workflow 조회 실패"`, `"호출 실패"`) + 영어(`"Not authorized"`, `"invalid path parameter"`) 혼재
- 동일 핸들러 내에서도 혼용. API 표면 일관성 부재

### M8. 동일 오류 메시지 반복 (디버깅 불가)

- `service/workflow/execute_workflow.go:74,80,86,114` — 전부 `"호출 실패"`
- 어떤 호출이 실패했는지 클라이언트/로그에서 구분 불가

### M9. import 그룹 정렬 미준수

- 거의 모든 .go 파일 — stdlib, external, internal이 혼재
- Go 표준: (1) stdlib (2) external (3) internal 순서

### M10. 패키지명 충돌 (workflow 서비스 내부)

- `service/workflow/create_workflow_version.go` — `import "github.com/example/zenflow/internal/workflow"`
- 현재 패키지도 `package workflow` → 이름 충돌. alias 미제공 시 컴파일 오류

### M11. 프론트엔드 workflow.version 필드 참조

- `frontend/src/pages/workflow-detail.tsx:60` — `workflow.version` 접근
- Workflow 모델에 `version` 필드 없음 → `undefined` 렌더링

### M12. Resume 버튼이 Activate와 동일 동작

- `frontend/src/pages/workflow-detail.tsx:61-62` — 둘 다 `activateWorkflowMutation.mutate({})`
- 상태별 조건부 렌더링 없음. 모든 버튼 항상 표시

### M13. React key={index}

- `frontend/src/pages/workflows.tsx:20` — 배열 index를 key로 사용
- 리스트 변경 시 컴포넌트 상태 혼란

### M14. Tailwind CSS 미완성

- `frontend/package.json` — tailwindcss, postcss, autoprefixer가 devDependencies에 있음
- `tailwind.config.js`, `postcss.config.js`, CSS @tailwind 디렉티브 전부 없음 → 데드 의존성

### M15. smoke.hurl source_workflow_id 하드코딩

- `tests/smoke.hurl:63` — `"source_workflow_id": 1` 고정값. 캡처 변수 `{{workflow_id}}` 미사용
- DB 시퀀스가 1이 아니면 테스트 실패

### M16. smoke.hurl 테스트 커버리지 부족

- ExecuteWorkflow, ExecuteWithReport, ArchiveWorkflow, DeleteSchedule, GetExecutionReport 미테스트
- paused → active 재전환 미검증

### M17. 불필요한 재조회 패턴

- `service/workflow/activate_workflow.go:75-83` — UpdateStatus 후 FindByID 재조회
- `service/template/clone_template.go:74-83` — IncrementCloneCount 후 재조회
- archive_workflow, pause_workflow 동일 패턴. 매번 불필요한 DB 라운드트립

### M18. OpenAPI 스키마 required 미정의

- `internal/api/types.gen.go` — 모든 모델 필드가 포인터 타입(`*string`, `*int64`, `*time.Time`)
- OpenAPI 컴포넌트 스키마에 `required` 미지정. maxLength/minLength 검증 태그도 미생성

---

## 분류: 1층 codegen 수정 / 2층 아키텍처 설계 / 3층 스펙 보완

### 1층: codegen 템플릿 수정 (패턴 고정, 기계적 수정 가능)

| 이슈 | 수정 방향 |
|------|-----------|
| C1. JWT 하드코딩 | 시크릿 미설정 시 fatal. Issue/Refresh/Verify 통일 |
| C7. go.mod 버전 | 실제 Go 버전으로 생성 |
| C8. defer 미실행 | signal.NotifyContext + srv.Shutdown 패턴 생성 |
| H5. context.Background() | 핸들러 컨텍스트를 model까지 전파 |
| H6. 페이지네이션 누락 | 모든 List에 LIMIT 적용 |
| H7. CORS | gin-contrib/cors 미들웨어 자동 생성 |
| H8. HTTP 상태 미검사 | fetch wrapper에 status 검사 포함 |
| H9. ActionCount | spec 변수 바인딩 재검토 (actions 결과 사용) |
| H10. 조회 결과 폐기 | 사용되지 않는 @get 결과 경고 또는 제거 |
| H11. SELECT * | 명시적 컬럼 목록 생성 |
| H12. ORDER BY 누락 | List 쿼리에 기본 정렬 추가 |
| H13. 모델 중복 | main.go에서 모델 1회 생성 후 공유 |
| H14. 중복 이메일 | unique constraint 위반 감지 → 409 |
| M7. 혼합 언어 | 에러 메시지 언어 통일 (설정 기반) |
| M8. 동일 오류 메시지 | 호출 지점별 고유 메시지 생성 |
| M9. import 정렬 | goimports 표준 그룹핑 적용 |
| M10. 패키지 충돌 | 동명 패키지 import 시 alias 자동 생성 |
| M11. 없는 필드 참조 | OpenAPI 스키마 기반 필드 매핑 검증 |
| M13. key={index} | item.id를 key로 사용 |
| M15. hurl 하드코딩 | 캡처 변수 활용 |
| M17. 불필요한 재조회 | UpdateStatus → RETURNING 활용 또는 인메모리 패치 |
| M18. required 미정의 | OpenAPI 스키마에 required 반영 → 비포인터 타입 생성 |

### 2층: 아키텍처 설계 변경 (TODO007 D1-D4)

codegen 템플릿 수정만으로는 해결 불가. SSaC 문법 확장, crosscheck 규칙 신설, 또는 codegen 아키텍처 결정이 필요.

| 이슈 | 핵심 문제 | 수정 방향 |
|------|-----------|-----------|
| C5. Rego 인증 미검증 (D1) | bearerAuth endpoint의 Rego rule이 claims 미참조 | crosscheck 규칙 추가: `Policy ↔ OpenAPI` claims 참조 검증 |
| C3. Publish/Deliver 순서 (D2) | `@publish`/`@call`의 트랜잭션 경계 개념 부재 | SSaC `@after-commit` 블록 또는 codegen 암묵 규칙 |
| H1. ResourceID 오류 (D3) | create/list에서 ResourceID에 넣을 값 규칙 없음 | crosscheck WARNING + SSaC `ResourceID: _` 생략 문법 |
| C2. 프론트엔드 인증 (D4) | 토큰 저장/주입 흐름 자체가 codegen 범위 밖 | STML codegen 확장 또는 fullend.yaml frontend.auth 설정 |

### 3층: 설계(specs) 레벨 보완

fullend를 아무리 고쳐도 스펙이 안 바뀌면 같은 코드가 나오는 것들.

| 이슈 | 보완 방향 |
|------|-----------|
| C4. Role 검증 + 공개 조직생성 | Rego에 register/createOrg 인가 규칙 추가. roles enum 정의 |
| C6. 크레딧 Race Condition | SSaC에 @lock 또는 atomic update 패턴 정의 |
| H2. SQL 인젝션 | BuildSelectQuery 내부에서 화이트리스트 재검증 강제 |
| H3. 이메일 열거 | 로그인 응답 균등화 (더미 해시 비교) |
| H4. VerifyToken 제로값 | 필수 클레임 누락 시 에러 반환 |
| M3. DISABLE_STATE_CHECK | dev-only 빌드 태그로 제한 |
| M5. NextRun 하드코딩 | cron 파서로 실제 다음 실행 시점 계산 |
| M16. 테스트 커버리지 | Execute, Archive, DeleteSchedule 등 시나리오 추가 |
