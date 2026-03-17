# Phase049: 런칭 준비 — CI, CHANGELOG, Getting Started

## 목표

GitHub 공개 전 HIGH 우선순위 3건을 완료한다.

## 원칙

- 사실만 적는다. 과장, 미래 약속, 없는 기능 언급 금지
- 증명 불가능한 내용은 포기한다
- 데모 프로젝트에 baseline 에러가 있으면 있는 그대로 보여준다

## 변경 파일 목록

### 1. GitHub Actions CI

- **파일**: `.github/workflows/test.yml`
- **내용**: `go test ./...` — 현재 20개 패키지 전체 통과. 이것만 돌린다
- **사실**: Go 1.22+, Ubuntu latest. lint는 현재 설정 없으므로 넣지 않음
- **뱃지**: README.md 상단에 CI passing 뱃지 추가

### 2. CHANGELOG.md

- **파일**: `CHANGELOG.md`
- **원칙**: git log에서 추출 가능한 사실만 기재. 꾸며내지 않음
- **범위**: v0.1.0 (초기 커밋 2026-03-08) ~ v0.1.49 (현재)
- **형식**: Phase 단위 요약. 커밋 메시지 기반

쓸 수 있는 것 (git log에 증거 있음):
- Phase001~030: CLI 골격, 9개 SSOT 파서, 교차 검증, 코드젠
- Phase031~043: filefunc 어노테이션 전수 적용 (1482 Go 파일)
- Phase044: OpenAPI 입력 검증 태그 코드젠
- Phase045~048: 뮤테이션 테스트 FAIL 0건 달성 (114 케이스, 97.8% 통과율)

**쓸 수 없는 것** (증거 없음):
- ~~성능 벤치마크~~ — 측정한 적 없음
- ~~사용자 피드백~~ — 사용자 없음
- ~~프로덕션 검증~~ — 더미 프로젝트만 테스트

### 3. Getting Started 5분 가이드

- **파일**: README.md 상단
- **내용**: zenflow-try05를 예제로 `fullend validate` 실행
- **사실**: zenflow-try05는 baseline에 12 errors, 7 warnings가 있음. 이건 더미 프로젝트의 의도적 불완전함이 아니라 pkg/auth funcspec 미매칭 + OpenAPI 제약 누락. 있는 그대로 보여주되 "이 에러들이 무엇인지" 설명
- **포기 항목**: `fullend gen` 데모 — zenflow-try05가 validate 통과하지 못하므로 gen 실행 불가. gen 데모는 거짓이 됨

### 포기 항목 정리

| 항목 | 사유 |
|---|---|
| `fullend gen` 데모 | zenflow-try05 baseline 에러로 gen 실행 불가. 거짓 데모가 됨 |
| `fullend chain` 데모 | chain은 작동하지만 출력이 validate만큼 직관적이지 않음. 런칭 후 추가 |
| 성능 수치 | 측정한 적 없음 |
| 데모 GIF | MEDIUM 우선순위. 런칭 후 추가 가능 |

## 의존성

- GitHub Actions: repository에 `.github/workflows/` 디렉토리 생성
- README.md: 현재 존재 여부 확인 필요

## 검증 방법

1. `.github/workflows/test.yml` push 후 Actions 탭에서 초록 뱃지 확인
2. CHANGELOG.md의 모든 항목이 git log에서 검증 가능한지 확인
3. Getting Started의 명령어를 그대로 복사해서 실행했을 때 문서와 동일한 출력인지 확인
