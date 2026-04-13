# Phase017 — RuntimeBugFixing Report

2026-04-14 · docker + hurl 실측 기반 런타임 버그 수정

## 수정 요약

| ID | 제목 | 파일 | 결과 |
|----|------|------|------|
| 1-1 | zenflow ActivateWorkflow 403 | pkg/authz, pkg/generate/gogin | ✓ |
| 1-2 | OPA_POLICY_PATH 디렉토리 지원 | pkg/authz/load_policy_from_path.go, resolve_policy_path.go | ✓ |
| 1-3 | DSN 기본값 개선 | pkg/generate/gogin/main_template.go | ✓ |
| 1-5 | gigbridge seed CHECK 위반 | dummys/gigbridge/specs/db/users.sql | ✓ |

## 1-1 root cause (디버깅 여정)

**초기 추정**: OPA owners 정적 로드 (잘못된 진단)

**실제 원인**: `pkg/authz/check.go` 의 `opaInput.claims` 가 `{user_id, role}` 만 포함, `org_id` 누락:

```go
// Before
"claims": map[string]interface{}{"user_id": req.UserID, "role": req.Role}
```

Rego 정책 `data.owners.workflow[input.resource_id] == input.claims.org_id` 에서 `input.claims.org_id` 가 undefined → 비교 실패 → deny.

**수정**: `CheckRequest.Claims map[string]any` 필드 추가. 생성기가 `authz.ClaimsFromStruct(currentUser)` emit. struct tag `authz:"<jwt_key>"` 로 fullend.yaml claims 매핑 보존.

```go
// CurrentUser struct (generated)
type CurrentUser struct {
    ID    int64  `authz:"user_id"`   // Go 필드명 "ID" ≠ JWT key "user_id"
    OrgID int64  `authz:"org_id"`
    Role  string `authz:"role"`
}
```

## 1-2 OPA path 개선

이전: `OPA_POLICY_PATH` env 필수 + 파일 경로만 수용 (디렉토리는 `is a directory` 에러).

이후:
- env 가 디렉토리면 `*.rego` glob 로 전체 로드 + concat
- env 미지정 시 `./internal/authz`, `./authz`, `./policy` 순 자동 탐색 (첫 존재 디렉토리)
- `DISABLE_AUTHZ=1` 은 그대로

## 1-3 DSN 기본값

이전: 하드코딩 `postgres://localhost:5432/app?sslmode=disable`

이후:
```go
dsnDefault := os.Getenv("DATABASE_URL")
if dsnDefault == "" {
    dsnDefault = "postgres://localhost:5432/gigbridge?sslmode=disable"  // 모듈명 기반
}
```

## 1-5 gigbridge seed

이전: `INSERT ... VALUES (0, ..., 'system', ...)` → `CHECK (role IN ('client','freelancer','admin'))` 위반.

이후: `role='admin'`. X-78 해소 + gen 재가동 + nobody 시드 DB 적용 가능.

## 실측 검증 (docker pg + hurl)

| 시나리오 | 결과 |
|---------|------|
| gigbridge smoke | **12/12 통과** ✓ |
| zenflow smoke | 7/N — register/login/Create/Activate/AddAction/Pause 통과, **ListWorkflows 에서 500** |

zenflow ListWorkflows 500 원인: 생성된 `List(opts)` 함수가 WHERE `org_id = $1` 하드코딩인데 args 에 org_id 값 주입 누락. SSaC `Workflow.List({Query: query})` 가 OrgID 전달 안 함. 생성기의 org 스코프 자동 주입 설계 이슈 — **Phase019 후보**.

## filefunc

- Phase017 신규 7 파일 모두 F1/F2/A10 준수
- baseline 위반 37 유지

## 영향 범위

`pkg/authz.CheckRequest` 에 `Claims` 필드 추가 → **기존 호출자 하위 호환** 유지 (Claims=nil 시 UserID/Role fallback).

생성기 변경은 새로 `gen` 하는 프로젝트에만 영향. 기존 artifacts 는 재생성 필요.

## 다음 Phase

**Phase018 — DDLPipelineIntegration**: DDL 위상정렬 + `schema.sql` 통합 산출 + `DEFAULT N FK` auto nobody seed. Phase017 1-5 의 수동 seed 수정을 자동화로 대체.

**Phase019 (신규 후보) — ORM/Scope**: SSaC `@get` 의 org 스코프 자동 주입, List 함수 시그니처에 claim 전달 설계 재정비.
