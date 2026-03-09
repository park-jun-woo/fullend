✅ 완료

# Phase 029: @component 폐지 — @func만 허용

## 목표

`@sequence call`에서 `@component` 지원을 완전 제거하고 `@func`만 허용한다.

## 배경

- `@component`는 interface만 선언하고 구현체 주입이 codegen에서 보장되지 않아 런타임 nil panic 발생
- notification 같은 복잡 기능은 `@func` 하나로도 해결 안 됨 — SSOT 범위 밖
- `@func pkg.funcName` 패턴으로 외부 로직 호출을 통일

ssac 측 수정은 `수정지시서014`로 별도 요청 완료. fullend 측 변경만 다룬다.

## 변경 항목

### A. dummy-lesson 스펙 정리

| 파일 | 변경 |
|------|------|
| `specs/dummy-lesson/service/enrollment/enroll_course.go` | `@sequence call` + `@component notification` 시퀀스 삭제 |
| `specs/dummy-lesson/model/notification.go` | 파일 삭제 |

### B. crosscheck — @component ERROR 추가

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/crosscheck.go` | SSaC 파싱 결과에서 `@component` 사용 감지 시 ERROR 발생: `"@component는 폐지되었습니다. @func pkg.funcName을 사용하세요"` |

### C. gluegen — @component 코드젠 제거

| 파일 | 변경 |
|------|------|
| `internal/gluegen/gluegen.go` | `collectComponents`, `collectComponentsForDomain` 함수 삭제. `transformSource`에서 component 치환 로직 삭제. components 파라미터 제거 |
| `internal/gluegen/server.go` | `generateServerStruct`에서 components 파라미터·Handler struct component 필드 생성 로직 삭제 |
| `internal/gluegen/domain.go` | `transformServiceFilesWithDomains`에서 components 관련 로직 삭제. Handler struct component 필드 생성 삭제 |

### D. manual-for-ai.md 문서 정리

| 위치 | 변경 |
|------|------|
| SSaC 문법 `@component` 태그 | 삭제 |
| 10 Sequence Types 표 `call` 행 | `@component or @func` → `@func` |
| `@component notification` 예시 | 삭제 |
| model/*.go Rules | `@component` 참조 설명 삭제. `// @dto` 설명만 유지 |
| SSOT Connection Map | `@component` 연결선 제거 |

### E. README.md

| 위치 | 변경 |
|------|------|
| Cross-Validation 목록 | `@component` 관련 항목 없음 (변경 불필요) |

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `specs/dummy-lesson/service/enrollment/enroll_course.go` | notification call 시퀀스 삭제 |
| `specs/dummy-lesson/model/notification.go` | **삭제** |
| `internal/crosscheck/crosscheck.go` | @component ERROR 검증 추가 |
| `internal/gluegen/gluegen.go` | component 관련 함수·파라미터 삭제 |
| `internal/gluegen/server.go` | components 파라미터·로직 삭제 |
| `internal/gluegen/domain.go` | components 관련 로직 삭제 |
| `artifacts/manual-for-ai.md` | @component 문서 전부 삭제 |

## 의존성

- ssac 수정지시서014 (ssac 측 @component 파서 제거) — 병렬 진행 가능. fullend가 crosscheck ERROR로 먼저 차단

## 검증 방법

```bash
# 1. fullend 빌드
go build ./cmd/fullend/

# 2. validate — @component 없는 스펙 통과 확인
fullend validate specs/dummy-lesson

# 3. gen + build + hurl 통과
fullend gen specs/dummy-lesson /tmp/dummy-lesson-out
cd /tmp/dummy-lesson-out/backend && go mod tidy && go build ./...

# 4. go test
go test ./internal/...
```

---

## 추가 수정 (스모크 테스트 중 발견)

### F. gluegen — transformSource receiver 중복 방지

| 파일 | 변경 |
|------|------|
| `internal/gluegen/gluegen.go` | `transformSource`에서 `\nfunc ` 뒤에 `(`가 이미 있으면(receiver 존재) 치환 건너뛰기. 기존 artifacts에 재실행 시 `(s *Server) (s *Server)` 중복 receiver 버그 수정 |

### G. authzgen — via 매핑 create 시 부모 테이블 직접 조회

| 파일 | 변경 |
|------|------|
| `internal/gluegen/authzgen.go` | `lookupOwner`에 `action` 파라미터 추가. `via` ownership 매핑에서 action이 `create`이면 자식 테이블 JOIN 대신 부모 테이블에서 직접 소유자 조회. 아직 자식 레코드가 없는 create 시 403 발생 버그 수정 |

### H. dummy-lesson model — @dto 추가

| 파일 | 변경 |
|------|------|
| `specs/dummy-lesson/model/session.go` | `IssueTokenResponse`, `HashPasswordResponse` DTO 추가. 교차 검증 경고 해소 |

### I. hurl-gen — filter 테스트 값 불일치 수정

| 파일 | 변경 |
|------|------|
| `internal/gluegen/hurl.go` | x-filter 쿼리 파라미터 값을 `test` → `test_string`으로 수정. `generateDummyValue`가 문자열 필드에 `test_string`을 생성하므로 filter 값도 일치시킴 |

### J. hurl-gen — state transition 순서 수정

| 파일 | 변경 |
|------|------|
| `internal/gluegen/hurl.go` | PUT 중 request body가 없는 것(PublishCourse 등 상태 전환)을 `transitionSteps`로 분리하여 create 직후, read 이전에 배치. 기존에는 update 그룹에 속해 read 뒤에 실행되어 `published = TRUE` 조건의 List가 빈 결과 반환 |

### K. x-include를 런타임 쿼리 파라미터에서 코드젠 메타데이터로 변경

`x-include`는 관련 테이블 JOIN을 코드젠 시 항상 적용하기 위한 메타데이터이므로, 런타임 `?include=` 쿼리 파라미터로 노출하면 안 된다.

| 파일 | 변경 |
|------|------|
| `internal/gluegen/queryopts.go` | `IncludeConfig` struct 삭제, `QueryOptsConfig`에서 `Include` 필드 삭제, `ParseQueryOpts`에서 `?include=` 파싱 로직 삭제 |
| `internal/gluegen/model_impl.go` | include 로직에서 `containsStr(opts.Includes, ...)` 조건 제거 → include 항상 실행 |
| `internal/gluegen/hurl.go` | `?include=` 쿼리 파라미터 생성 제거 |
| `internal/gluegen/gluegen.go` | `buildQueryOptsConfig`에서 `x-include` → `Include:` 설정 생성 제거 (TODO) |
| SSaC 수정지시서015 | `QueryOpts`에서 `Includes []string` 필드 제거 (정리 차원)
