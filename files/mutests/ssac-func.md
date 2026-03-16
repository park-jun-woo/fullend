# Mutation Test — SSaC → Func

## CheckFuncs (SSaC @call → Func Spec)

### MUT-SSAC-FUNC-001: SSaC @call 패키지명 오타
- 대상: `specs/gigbridge/service/auth/login.ssac`
- 변경: `auth.IssueToken` → `Auth.IssueToken`
- 기대: ERROR — funcspec package "auth"와 불일치
- 결과: PASS — Func↔SSaC 패키지 불일치 검출

### MUT-SSAC-FUNC-002: @call input 타입 불일치
- 대상: `specs/gigbridge/service/billing/charge.ssac`
- 변경: @call input에 int 타입을 string 파라미터 위치에 전달
- 기대: ERROR — 타입 불일치
- 결과: 미실행

### MUT-SSAC-FUNC-003: @call func이 미참조 (func_coverage)
- 대상: func/ 디렉토리에 SSaC에서 @call하지 않는 함수 추가
- 변경: `func/billing/unused.go`에 `// @func Unused` 추가
- 기대: WARNING — custom func spec이 SSaC @call에서 미참조
- 결과: 미실행

## CheckAuthz (SSaC @auth → pkg/authz CheckRequest)

### MUT-AUTHZ-001: @auth 잘못된 input 필드
- 대상: `specs/gigbridge/service/gig/create_gig.ssac`
- 변경: @auth input에 `Department: "engineering"` 추가 (허용 필드: Action, Resource, UserID, Role, ResourceID)
- 기대: ERROR — CheckRequest에 "Department" 필드 없음
- 결과: 미실행
