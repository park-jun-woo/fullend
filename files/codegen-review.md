# fullend codegen 프로덕션 리뷰

대상: `artifacts/dummys/zenflow-try06/` (fullend gen 생성 코드)
일시: 2026-03-18
기준: 현업 프로덕션 SaaS 배포 가능 수준

---

## 요약

| 심각도 | 건수 | 주요 영역 |
|--------|------|-----------|
| CRITICAL | 6 | JWT 위조, 인증 미작동, Race Condition |
| HIGH | 9 | 테넌트 격리, SQL 인젝션, 컨텍스트 누락 |
| MEDIUM | 13 | 운영 안정성, 프론트엔드 품질, 관측성 |

---

## CRITICAL — 배포 불가

### C1. JWT 시크릿 하드코딩

- `backend/internal/auth/issue_token.go:26-29` — `JWT_SECRET` 미설정 시 `"secret"` 폴백
- `backend/internal/auth/refresh_token.go:33` — `"secret"` 리터럴 고정
- **영향**: 공격자가 임의 토큰 위조 가능. 전체 인증 체계 무력화
- **원인**: codegen이 시크릿 폴백 로직을 하드코딩으로 생성
- **수정**: 시크릿 미설정 시 서버 시작 거부 (fatal). refresh_token도 주입 방식으로 전환

### C2. 프론트엔드 토큰 미사용

- `frontend/src/api.ts` 전체
- **현상**: login 응답의 `access_token`을 저장하지 않고, 후속 요청에 Authorization 헤더 미부착
- **영향**: 인증 필요 API 전부 401. 프론트엔드 사실상 미작동
- **원인**: codegen이 토큰 저장/주입 코드를 생성하지 않음
- **수정**: api.ts에 토큰 저장소 + fetch wrapper에 Bearer 헤더 자동 부착 생성

### C3. 크레딧 차감 Race Condition

- `backend/internal/service/workflow/execute_workflow.go:72-93`
- `backend/internal/service/workflow/execute_with_report.go:73-94`
- **현상**: CheckCredits → DeductCredit 사이에 행 잠금 없음
- **영향**: 동시 요청 시 크레딧 이중 차감, 음수 잔액 발생
- **수정**: `SELECT ... FOR UPDATE` 또는 `UPDATE ... WHERE credits >= cost RETURNING credits` 원자적 차감

### C4. Queue Publish가 tx.Commit 이전

- `backend/internal/service/workflow/execute_workflow.go:101-109`
- `backend/internal/service/workflow/execute_with_report.go:108-116`
- **현상**: queue.Publish()가 tx.Commit() 보다 먼저 호출됨
- **영향**: Commit 실패 시 이미 발행된 메시지 회수 불가. 웹훅이 미존재 실행에 발송
- **수정**: Publish를 Commit 이후로 이동하거나 Transactional Outbox 패턴 적용

### C5. Webhook 동기 호출 (트랜잭션 내)

- `backend/internal/service/workflow/execute_workflow.go:111-114`
- `backend/internal/service/workflow/execute_with_report.go:118-121`
- **현상**: 트랜잭션 내에서 외부 HTTP 호출 (webhook.Deliver)
- **영향**: 외부 endpoint 지연 시 DB 트랜잭션 장시간 점유, 커넥션 고갈
- **수정**: 웹훅은 Commit 이후 비동기 전송. 실패 시 재시도 큐 활용

### C6. Role 검증 없는 회원가입

- `backend/internal/service/auth/register.go:15`
- **현상**: 요청 바디의 `Role` 필드를 검증 없이 그대로 저장. `binding:"required,max=50"`만 존재
- **영향**: `"admin"`, `"super_admin"` 등 임의 역할로 가입 → 권한 상승
- **수정**: codegen이 fullend.yaml의 roles 목록으로 enum 검증 코드 생성

---

## HIGH — 보안/데이터 무결성 위험

### H1. 테넌트 격리 누락 (GetTemplate)

- `backend/internal/service/template/get_template.go:18`
- authz.Check 없이 template ID만으로 조회. 타 조직 템플릿 열람 가능

### H2. authz에 잘못된 ResourceID

- `backend/internal/service/workflow/create_workflow.go:32`
- `backend/internal/service/workflow/list_workflows.go:14`
- `ResourceID: currentUser.ID` — 워크플로우 ID가 아닌 유저 ID로 인가 체크

### H3. DB 조회 후 인가 체크 (순서 역전)

- `backend/internal/service/log/get_execution_report.go:22-33`
- `backend/internal/service/template/clone_template.go:29-43`
- `backend/internal/service/template/publish_template.go:36-49`
- `backend/internal/service/action/add_action.go:42-56`
- FindByID 먼저 실행 → authz.Check 나중. 리소스 존재 여부 타이밍 공격 가능

### H4. QueryOpts SQL 인젝션

- `backend/internal/model/queryopts.go:131`
- 필터 컬럼명이 `fmt.Sprintf("%s = $%d", col, ...)` 로 직접 삽입
- ParseQueryOpts의 화이트리스트 검증이 불완전하여 우회 가능성 존재

### H5. 이메일 열거 공격

- `backend/internal/service/auth/login.go:22-30`
- 유저 미존재 시 즉시 404, 존재 시 bcrypt 검증 → 응답 시간 차이로 이메일 유효성 확인
- **수정**: 미존재 시에도 더미 해시 비교 실행하여 타이밍 균등화

### H6. 페이지네이션 없는 List 엔드포인트

- `backend/internal/service/log/list_execution_logs.go:38`
- `backend/internal/service/workflow/list_workflows.go:19`
- `backend/internal/service/webhook/list_webhooks.go:19`
- 행 수 제한 없이 전체 반환. 대량 데이터 시 OOM

### H7. context.Background() 일괄 사용

- `backend/internal/model/workflow.go:44,52,60,90` 외 전체 model 파일
- 요청 컨텍스트 미전파. 쿼리 취소/타임아웃 불가, 클라이언트 끊겨도 쿼리 계속 실행

### H8. CORS 미설정

- `backend/internal/service/server.go:28-29`
- gin.Default()만 사용. CORS 미들웨어 없음. 프론트엔드 크로스 오리진 요청 차단

### H9. 프론트엔드 HTTP 상태 미검사

- `frontend/src/api.ts:17,34,51` 외 전체 API 함수
- `res.json()` 직접 호출. 500 에러 응답도 성공 데이터로 처리

---

## MEDIUM — 운영 안정성/품질

### M1. Graceful Shutdown 없음

- `backend/cmd/main.go:110-114`
- queue.Start() 고루틴 + log.Fatal(r.Run()). 시그널 핸들링 없음. 종료 시 진행 중 작업 유실

### M2. DB 커넥션 풀 미설정

- `backend/cmd/main.go:38`
- SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime 미호출. 고부하 시 커넥션 고갈

### M3. 상태 머신 디버그 우회

- `backend/internal/states/workflowstate/workflowstate.go:40-41`
- `DISABLE_STATE_CHECK=1` 환경변수로 전체 상태 검증 무력화. 프로덕션 사고 위험

### M4. Health/Ready 엔드포인트 없음

- `backend/internal/service/server.go`
- K8s liveness/readiness probe 불가. 오토스케일링/롤링 배포 지원 불가

### M5. 구조화된 로깅 없음

- 전체 백엔드
- request ID, trace ID 없음. `log.Printf()` 만 사용. 분산 환경에서 디버깅 불가

### M6. Rate Limiting 없음

- 전체 엔드포인트
- 로그인/가입 brute force 무방비

### M7. 멱등성 키 없음

- `backend/internal/service/workflow/execute_workflow.go`
- `backend/internal/service/workflow/execute_with_report.go`
- 네트워크 재시도 시 이중 실행 + 이중 과금

### M8. CloneTemplate 순서 오류

- `backend/internal/service/template/clone_template.go:68`
- IncrementCloneCount가 워크플로우 생성 성공 여부와 무관하게 실행

### M9. os.Setenv로 JWT 시크릿 설정

- `backend/cmd/main.go:48`
- 프로세스 환경변수에 시크릿 노출. `/proc/PID/environ` 으로 탈취 가능

### M10. go 1.25.0

- `backend/go.mod:3`
- 존재하지 않는 Go 버전. 빌드 불가

### M11. React key={index}

- `frontend/src/pages/workflows.tsx:21`
- 리스트 key에 배열 index 사용. 리렌더링 시 컴포넌트 상태 혼란

### M12. TypeScript any 남발

- `frontend/src/pages/workflow-detail.tsx:16,23,37,44`
- API 응답에 타입 안전성 없음. 런타임 에러 원인

### M13. CSP 헤더 없음

- `frontend/index.html`, `backend/internal/service/server.go`
- Content-Security-Policy 미설정. XSS 실행 차단 레이어 부재

---

## 분류: codegen 수정 vs 설계 보완

### codegen(fullend)이 고쳐야 할 것

| 이슈 | 수정 방향 |
|------|-----------|
| C1. JWT 하드코딩 | 시크릿 미설정 시 fatal 종료 생성 |
| C2. 프론트엔드 토큰 | 토큰 저장소 + fetch wrapper 생성 |
| C4. Publish 순서 | Publish를 Commit 이후로 배치 |
| C5. Webhook 동기 | Commit 이후 비동기 호출로 변경 |
| C6. Role 검증 | fullend.yaml roles로 enum 검증 생성 |
| H2. ResourceID 오류 | create/list에 적합한 ResourceID 로직 |
| H3. 인가 순서 역전 | authz.Check를 DB 조회 이전으로 |
| H6. 페이지네이션 누락 | 모든 List에 ParseQueryOpts 적용 |
| H7. context.Background() | 핸들러 컨텍스트를 model까지 전파 |
| H8. CORS | gin-contrib/cors 미들웨어 자동 생성 |
| H9. HTTP 상태 미검사 | fetch wrapper에 status 검사 포함 |
| M1. Graceful Shutdown | signal.NotifyContext + srv.Shutdown 생성 |
| M2. 커넥션 풀 | 환경변수 기반 풀 설정 코드 생성 |
| M4. Health 엔드포인트 | /health, /ready 자동 생성 |
| M10. go.mod 버전 | 실제 Go 버전으로 생성 |
| M11. key={index} | item.id를 key로 사용 |
| M12. any 타입 | OpenAPI 스키마에서 TS 타입 생성 |

### 설계(specs) 레벨에서 보완할 것

| 이슈 | 보완 방향 |
|------|-----------|
| C3. 크레딧 Race Condition | SSaC에 @lock 또는 atomic update 패턴 정의 |
| H1. 테넌트 격리 | Rego 정책에 GetTemplate 인가 규칙 추가 |
| H4. SQL 인젝션 | QueryOpts 화이트리스트 강제 적용 구조 |
| H5. 이메일 열거 | 로그인 응답 균등화 정책 |
| M3. 상태 머신 우회 | DISABLE_STATE_CHECK를 dev-only 빌드 태그로 제한 |
| M6. Rate Limiting | fullend.yaml에 rate limit 설정 → 미들웨어 생성 |
| M7. 멱등성 키 | OpenAPI x-idempotency 확장 → 미들웨어 생성 |
| M9. 시크릿 관리 | DI 패턴으로 시크릿 전달. os.Setenv 폐기 |
| M13. CSP | fullend.yaml에 보안 헤더 설정 → 미들웨어 생성 |
