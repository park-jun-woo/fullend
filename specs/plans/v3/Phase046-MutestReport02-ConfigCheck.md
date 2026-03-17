# Phase046: Config 교차 검증 강화 ✅ 완료

## 목표

mutest-report02 FAIL 중 config-check / func-check 도메인을 해결한다. 확정 수정 2건, 보류 1건.

## 사전 재확인 결과

| ID | 계획 시 판단 | 재확인 결과 |
|---|---|---|
| MUT-SSAC-CONFIG-001 | 기존 로직 있음 → agent 오류 의심 | **PASS 확인** — agent가 `currentUser.ID` 삭제했으나 zenflow SSaC는 `currentUser.OrgID`만 참조. `OrgID` 삭제 시 정상 검출. **제외** |
| MUT-SSAC-037 | 기존 로직 있음 → agent 오류 의심 | **진짜 FAIL** — `strcase.ToGoCamel("IssueToken")` = `"issueToken"` (lowerCamelCase). key `"auth.issueToken"` → `jwtBuiltinFuncs` 매칭 → 모든 검증 스킵. OrgId 변이 검출 불가 |

## 변경 파일 목록

### 1. MUT-POLICY-CONFIG-001 — Rego role 값 대소문자 미검출

- **chain**: `CheckRoles` ← `rules.go`
- **현상**: `CheckRoles`는 `fullend.yaml auth.roles` 목록과 Rego role 값을 비교. zenflow에는 roles 섹션이 없어서 `CheckRoles` 자체가 실행되지 않음. Rego의 `input.claims.role == "Admin"` (대문자) 오타를 검출할 수단 없음
- **수정**: DDL `CHECK (role IN ('admin', 'member'))` 제약에서 유효한 role 값을 추출하여 Rego role 값과 교차 검증. roles 섹션 없이도 DDL CHECK 기반으로 검증 가능
- **파일**: `internal/crosscheck/` 신규 파일 (check_rego_role_ddl.go)
- **선행 확인**: `SymbolTable`에 DDL CHECK 제약 정보가 있는지 확인. 없으면 DDL 파서 확장 필요
- **테스트**: Rego role "Admin" + DDL CHECK "admin" → ERROR 검출

### ~~2. MUT-POLICY-CONFIG-002~~ — 제외 (검증 범위 밖)

claims JWT key(`user_id`)는 사람이 자유롭게 정하는 이름. 코드젠이 발급(`generate_issue_token.go:23`)과 파싱(`claim_extract_line.go:17`) 양쪽 모두 `def.Key`를 사용하므로, `user_id` → `userId` 변경 시 자기 일관적. `go build` 깨지지 않음. SSOT 내 교차 검증 근거 없음

### 3. MUT-SSAC-037 — jwtBuiltinFuncs가 @call 입력 검증을 스킵

- **chain**: `checkSingleCall` → `jwtBuiltinFuncs[key]` → `return nil`
- **현상**: `parseCallKey("auth.IssueToken")` → `strcase.ToGoCamel("IssueToken")` = `"issueToken"` → key `"auth.issueToken"` → `jwtBuiltinFuncs["auth.issueToken"]` = true → **모든 검증 스킵**. jwtBuiltinFuncs는 코드젠 대상 함수라 funcspec이 없어서 스킵하는데, 이 때 input key 검증까지 스킵됨
- **수정**: jwtBuiltinFuncs 스킵 시에도 input key의 기본 검증(claims 필드와의 매칭)을 수행. 또는 jwt 빌트인 함수의 Request 필드를 하드코딩하여 specMap에 등록 (IssueTokenRequest: ID, Email, Role, OrgID 등)
- **대안**: `checkServiceFuncCalls`에서 jwtBuiltinFuncs 대상 @call의 input key가 claims 키와 매칭되는지 별도 검증
- **파일**: `internal/crosscheck/check_single_call.go` 또는 신규 파일
- **테스트**: `@call auth.IssueToken({..., OrgId: user.OrgID})` 시 "OrgId"가 claims 키에 없음 ERROR 검출

### 제외 항목

- **MUT-SSAC-CONFIG-001**: 기존 코드 정상 동작 확인. agent가 잘못된 claim 삭제
- **MUT-SCENARIO-002**: 빈 scenario WARN은 의도된 동작 (선택적 SSOT)

## 의존성

- 수정 1: `DDLTable.CheckEnums` (`map[string][]string`) 이미 존재 — DDL 파서 확장 불필요
- 수정 1: `collectRegoRoles` 이미 존재 — Rego role 값 수집 재사용
- 수정 3: `CrossValidateInput.Claims` (`map[string]ClaimDef`) 이미 존재 — claims 키 목록으로 jwt 빌트인 함수 input key 검증
- 추가 외부 패키지 없음

## 검증 방법

1. `go test ./internal/crosscheck/...` 통과
2. `go test ./...` 통과
3. zenflow-try05 대상 2건 뮤테이션 재실행으로 PASS 전환 확인 (MUT-POLICY-CONFIG-001, MUT-SSAC-037)
