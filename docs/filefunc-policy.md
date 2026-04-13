# filefunc 규약 정책 — fullend 프로젝트 결정 (Phase015 Part C)

2026-04-14 · v5 로드맵 마감 시점 결정

## 배경

filefunc (`~/.clari/repos/filefunc`) 은 Go 소스에 대해 다음 규약을 강제:

- **F1**: 파일당 1 func
- **F2**: 파일당 1 type
- **A1**: 모든 func/type 에 `//ff:func` / `//ff:type` 어노테이션
- **A2**: codebook.yaml 에 feature/type/topic 선등록
- **A3**: `//ff:what` 한 줄 요약
- **A10~A13**: `control=selection/iteration/sequence` 가 실제 제어 흐름과 일치
- **Q1**: 중첩 depth ≤ 2
- **Q4**: range/loop 본문 ≤ 10 줄

fullend 프로젝트의 Phase001~018 진행 경험으로 규약의 **비용·편익 관찰**이 축적됨. 본 문서는 최종 정책 결정.

## Phase010 에서 드러난 비용

`Decide*` 패턴 수렴 시 응집 단위 강제 분해:

| 논리 단위 | F1/F2 강제 결과 |
|---------|---------------|
| `Pattern` enum + `MethodFacts` struct + `NewMethodFacts` + `DecideMethodPattern` | **4 파일로 분할** (같은 주제) |
| Decide* 3개 × 4 atoms = **12 파일** 생산 (Phase010) |

Go stdlib 관례와 충돌:
- `net/http/request.go` — 수십 개 func/type 공존
- `encoding/json/encode.go` — 약 2000줄 단일 파일
- fullend pkg/generate/gogin 에서 동일 수준 응집을 강제 분해 → 파일 ~470개

테스트 파일 역설:
- `_test.go` 에서 F1 적용 시 table-driven 테스트 (1 func 안 여러 subtest) 와 충돌
- baseline 위반에 test 파일들 다수 포함 — 실제론 묵인 중

## Phase010~018 에서 선택한 대응 (사실상의 정책)

내 파일들은 아래 원칙으로 F1/F2 준수:

1. **Production 파일**: 엄격 1 func/type/file. 응집 atom 4개면 4 파일.
2. **Test 파일**: F1 위반 허용 (baseline 과 일치). `_test.go` 는 table-driven 관례 유지.
3. **제어 흐름 annotation (A10~A13)**: 엄격 준수. 리팩토링 시 annotation 도 갱신.
4. **codebook 선등록 (A2)**: 새 feature/type/topic 은 먼저 codebook.yaml 에 추가.

이 운영으로 **baseline 위반 37 유지** (Phase011~018 내내). 내 신규 파일 위반은 매 phase 0.

## 선택안 비교

| 안 | 장점 | 단점 | 평가 |
|----|------|------|------|
| **A. F1/F2 엄격 유지** | 명확한 규칙, whyso/filefunc 인덱싱 용이 | 응집 atom 강제 분해, 파일 수 폭증 | **현행** |
| **B. `//ff:group` 도입** | 응집 단위 한 파일 허용 | filefunc 도구 확장 필요, 그룹 경계 모호 | **중기 연구 대상** |
| **C. 테스트 파일 F1 면제** | 관례 일치, 사용자 비용 낮음 | F1 을 부분 면제하면 일관성 약화 | **실질 이미 반영** |
| **D. primary symbol 규칙** | 가장 자연스러움 (1 파일 1 주요 심볼 + helper) | primary 판정 모호, tooling 구현 복잡 | **장기 목표** |

## 결정

**본 프로젝트 (fullend v5 종료 시점) 는 안 A (F1/F2 엄격 유지) + 사실상 C 의 혼합**을 공식화:

1. **Production `.go`**: F1/F2 엄격. 응집 atom 은 파일로 분할, 디렉토리 구조로 응집 표현.
2. **`_test.go`**: F1/F2 면제 (사실 baseline 이 그러함). test-helper 는 별도 파일 권장이나 강제 안 함.
3. **Annotation (A1/A2/A3/A10~A13)**: 예외 없이 엄격.
4. **Q1/Q4 복잡도 상한**: 엄격.

### 재검토 조건

다음 중 하나라도 발생하면 정책 재평가:

- filefunc 리포에 `//ff:group` 또는 primary-symbol 규칙 구현 반영
- fullend 파일 수가 pkg 1000개 넘음 (현재 1200)
- 신규 phase 에서 F1/F2 비용이 기능 진전을 지연시키는 사례 발생

### 재검토 주체

fullend 유지보수자 + filefunc 도구 개발자 간 협의.

## 관련 이슈 후보 (filefunc 리포에 제출 권장)

- "test 파일 F1/F2 공식 면제" — 현 baseline 이 이미 반영 중이나 명문화 필요
- "`//ff:group` primitive" — 응집 atom 한 파일 허용
- "primary symbol 규칙 연구" — 1 파일 1 주요 심볼 + helper/type 허용

## 운영 지침

### 새 파일 작성 시 체크리스트

1. `//ff:func` / `//ff:type` + `//ff:what` 한 줄
2. codebook.yaml 에 사용할 feature/type/topic 등록 확인
3. 파일당 func 1개, type 1개 (production 한정)
4. 제어 흐름 annotation 실제 구조와 일치
5. 중첩 ≤ 2, range body ≤ 10줄

### `filefunc validate` 주기적 실행

- 각 phase 커밋 전 `filefunc validate` 으로 baseline 증감 확인
- 신규 파일 위반 0 유지
- baseline 위반은 기존 파일 유래이며 장기적 감소 목표 (별도 cleanup phase)

## baseline 현황 (2026-04-14, Phase015 완료 시점)

- **baseline 37** (일관 유지)
- 주된 baseline 발생원:
  - 일부 test 파일 F1 위반 (table-driven 관례)
  - 기존 `internal/` 파일들 (Phase014 삭제 보류로 정리 미완)
  - 일부 codebook 누락 (historical)

baseline 0 달성은 **filefunc 정책 재검토 + 장기 cleanup** 조합 필요.
