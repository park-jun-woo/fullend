# TODO004: dummy-GigBridge 테스트 중 발견된 fullend 미비점

## 1. @auth → pkg/authz 기반 @call func 방식 전환

- **현재**: fullend authzgen이 고정 `authz.Input{ID interface{}}` 자동생성 → SSaC inputs와 불일치
- **방안**: `pkg/authz`를 기본 제공. fullend.yaml에 `authz:` 설정 없으면 `pkg/authz` 기본값, 있으면 사용자 지정 패키지
- **fullend 작업**:
  - `pkg/authz/` 신규 작성 (기본 OPA Rego 구현)
  - `fullend.yaml` 파서: `authz.package` 필드 추가
  - `authzgen` 자동생성 제거 → `@call func`과 동일한 복사 방식
  - crosscheck: `@auth` inputs → authz CheckRequest 필드 매칭 검증
- **SSaC 수정지시서**: `~/.clari/repos/ssac/files/수정지시서v2/수정지시서022.md`

## 2. @call func 타입 불일치 검증 (SSaC 영역)

- **현재**: SSaC validator가 `@call` inputs의 필드 이름만 검증, 타입은 비교 안 함
- **문제**: DDL `INTEGER` → sqlc `int32`, func Request `int` → 타입 변환 없이 대입 → 컴파일 에러
- **방안**: SSaC validator에서 변수의 DDL 출처 타입과 func Request 필드 타입을 비교 → 불일치 시 ERROR
- **SSaC 수정지시서**: `~/.clari/repos/ssac/files/수정지시서v2/수정지시서022.md`에 포함

## 3. config.* → pkg/config 기본 조회 함수

- **현재**: SSaC `config.SMTPHost` → 코드젠이 `config.SMTPHost` 그대로 출력 → `undefined: config`
- **방안**: `pkg/config`에 환경변수 조회 함수 제공. SSaC codegen이 `config.Key` → `config.Get("KEY")` 변환
- **대상**: fullend `pkg/config/` 신규 + ssac generator 변경
- **SSaC 수정지시서**: `~/.clari/repos/ssac/files/수정지시서v2/수정지시서022.md`에 포함

## 4. @call/@post 미사용 변수 _ 처리

- **현재**: 결과 변수가 이후 @response에서 미참조 → Go `declared and not used` 에러
- **방안**: generator가 미참조 변수를 `_`로 생성
- **SSaC 수정지시서**: `~/.clari/repos/ssac/files/수정지시서v2/수정지시서022.md`에 포함
