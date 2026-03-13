# TODO006: dummy-GigBridge 개발 결과 및 발견 버그

## 개발 개요

- **프로젝트**: GigBridge (Freelance Escrow Matching Platform)
- **소요 시간**: 약 22분 38초 (01:51:29 ~ 02:14:07)
- **SSOT 파일**: 11개 SSOT + Func Spec (총 25+ 파일)

## 테스트 결과

| 테스트 | 결과 | 비고 |
|---|---|---|
| `fullend validate` | PASS | 0 errors, 27 warnings |
| `fullend gen` | PASS | 12/12 codegen 통과 |
| `go build` | PASS | |
| `smoke.hurl` | PASS | 12/12 requests (DISABLE_AUTHZ=1) |
| `invariant.hurl` | FAIL | 403 기대 → 200 반환 |

## 발견 버그

### BUG-1: authz-gen — OPA data.owners 미로딩 (CRITICAL)

- **위치**: `internal/gluegen/authz_gen.go` (추정)
- **현상**: `@ownership` 어노테이션에 정의된 소유권 데이터(`data.owners`)가 OPA 평가 전에 DB에서 조회·로딩되지 않음
- **영향**: `Check()` 함수가 `input`만 전달하고 `data.owners`를 비워둬서 모든 소유권 기반 allow 룰이 항상 실패
- **재현**: DISABLE_AUTHZ 없이 서버 기동 → PublishGig 호출 → 소유자여도 403 반환
- **기대 동작**: `Check()` 호출 시 `@ownership gig: gigs.client_id` 기반으로 DB 쿼리 실행 → `data.owners.gig[resource_id]` = owner_user_id 매핑 후 OPA 평가
- **추가 확인**: OPA input 키도 `resource_owner_id`로 설정되어 있으나, Rego 정책은 `input.resource_id`를 참조 — 키 이름 불일치

### BUG-2: hurl-gen — 중첩 객체 토큰 캡처 (MINOR, 우회 완료)

- **위치**: `internal/gluegen/hurl_gen.go` (추정)
- **현상**: Login 응답이 `{"token": {"access_token": "..."}}` 구조일 때, smoke.hurl이 `jsonpath "$.token"`으로 캡처하여 객체를 Bearer 토큰으로 사용 시도 → Hurl 렌더링 에러
- **우회**: Login SSaC `@response`를 `{ token: token }` → `token`으로 변경하여 플랫 응답으로 전환
- **근본 원인**: hurl-gen이 중첩 응답 구조에서 토큰 필드 경로를 자동 추출하지 못함

### BUG-3: pkg/auth — IssueTokenResponse JSON 태그 누락 (MINOR)

- **위치**: `pkg/auth/issue_token.go:19`
- **현상**: `AccessToken string` 필드에 JSON 태그 없음 → Go 기본 직렬화로 `{"AccessToken": "..."}` 출력 (PascalCase)
- **영향**: OpenAPI snake_case 규칙(`access_token`)과 불일치
- **제안**: `AccessToken string \`json:"access_token"\`` 태그 추가

## SSOT 작성 중 학습한 규칙

1. **SSaC 파일당 1 함수**: `service/` 직접 배치 불가, 도메인 서브폴더 필수 (예: `service/gig/create_gig.ssac`)
2. **Go 예약어 컬럼 금지**: DDL `type` → `tx_type` 등으로 변경 필요
3. **Page[T] 응답**: OpenAPI에 `items` + `total` 필드만 선언 (limit/offset은 프레임워크 자동 처리)
4. **OPA v1 문법**: `allow if {` (if 키워드 필수)
5. **SSaC import 필수**: `@call pkg.Func` 사용 시 해당 패키지 import 선언 필요
6. **model/ 디렉토리 필수**: @dto 타입 없어도 빈 `model/model.go` (package 선언만) 필요
7. **@put 후 재조회**: `@put Model.Update(...)` 후 `@response`에 사용하려면 `@get`으로 재조회 필요 (WARN 방지)
