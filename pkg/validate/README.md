# pkg/validate

단일 SSOT 자체 정합성 검증.

## 심볼 테이블

| 패키지 | 설명 |
|--------|------|
| `symbol/` | 이름→타입 조회 허브. DDL/OpenAPI/Model/SQLc에서 구축. 모든 validator가 공유 |

## 개별 검증

| 패키지 | 대상 | 설명 |
|--------|------|------|
| `ssac/` | `[]ServiceFunc` | 변수 흐름, 필수 필드, stale response, Go 예약어 등 |
| `stml/` | `[]PageSpec` | fetch/action 바인딩, 파라미터, 컴포넌트 참조 등 |

## Toulmin 매핑

```
claim   = 검증 대상 (ServiceFunc, PageSpec)
ground  = 심볼 테이블 (SymbolTable)
backing = 검증 기준/설정 (nil 또는 고정 설정)

warrant  = 기본 규칙 ("변수가 선언 후 사용되어야 한다")
rebuttal = 예외 ("currentUser는 암묵적 선언")
```
