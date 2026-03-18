# Phase007 — E2E 검증: zenflow-try06 + 신규 dummy

## 목표

toulmin 전환 후 전체 파이프라인을 E2E로 검증한다. zenflow-try06 기존 프로젝트 + 복잡 분기 시나리오로 실증.

## 검증 항목

### 1. zenflow-try06 전환 검증

```bash
# specs 디렉토리에 authz/*.go 존재, policy/*.rego 부재 확인
ls specs/dummys/zenflow-try06/authz/
ls specs/dummys/zenflow-try06/policy/  # should not exist

# validate
fullend validate specs/dummys/zenflow-try06/

# gen
fullend gen specs/dummys/zenflow-try06/ artifacts/dummys/zenflow-try06/

# build
cd artifacts/dummys/zenflow-try06/backend && go build ./...

# .rego 파일 0개 확인
find artifacts/dummys/zenflow-try06/ -name "*.rego" | wc -l  # 0

# OPA import 0개 확인
grep -r "open-policy-agent" artifacts/dummys/zenflow-try06/  # 0

# hurl 스모크 테스트
hurl --test artifacts/dummys/zenflow-try06/tests/smoke.hurl
```

### 2. 인가 동작 검증

기존 Rego 규칙과 동일한 인가 판정 확인:

| 시나리오 | 기대 결과 |
|----------|----------|
| admin + 모든 API | 허용 |
| member + 자기 org의 workflow | 허용 |
| member + 타 org의 workflow | 거부 |
| 미인증 + bearerAuth endpoint | 거부 |
| 미인증 + public endpoint | 허용 |

### 3. defeats 그래프 검증

```bash
# YAML 그래프 정의 검증 (순환 탐지 + 참조 유효성)
toulmin graph specs/dummys/zenflow-try06/authz/authz.yaml --check

# 생성된 authz 코드에서 defeats 그래프 분석
cd artifacts/dummys/zenflow-try06/backend
toulmin graph internal/authz/graph_gen.go
```

defeats 관계가 specs의 YAML 정의와 일치하는지 확인.

### 4. h-Categoriser verdict 검증

`EvaluateTrace`로 판정 경로 추적:
- admin이지만 suspended → SuspendedBlock이 AdminAllowAll을 무력화 → DenyByDefault verdict > 0 → 거부
- owner이지만 suspended → SuspendedBlock이 OwnerAccess를 무력화 → 거부
- member + org 소속 + org visibility → OrgVisibility가 DenyByDefault를 무력화 → verdict <= 0 → 허용
- TraceEntry에서 각 규칙의 activated/role/qualifier 확인

### 5. codegen-review2.md 이슈 해소 확인

| 이슈 | 상태 |
|------|------|
| C5. Rego 8개 규칙 인증 미검증 | 해소: Go 함수에서 claims 참조가 코드에 명시 |
| D1. Rego claims 검증 누락 탐지 | 해소: AllowRule.UsesClaims로 crosscheck |
| D3. ResourceID 바인딩 | 해소: Go 함수가 직접 claim 필드 접근 |

## 의존성

- Phase001~006 전체 완료 필수

## 검증 방법

- 위 전체 항목 통과
- `go test ./...` 전체 통과
